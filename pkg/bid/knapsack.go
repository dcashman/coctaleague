package bid

// The fundamental unit for knapsack calculation: we need to be able to judge item cost (classicall weight) and benefit (value).
type Item interface {
	Value() int
	Weight() int
}

// ItemProgression: basically a linked list connecting each item we collect in our knapsack so we can retrieve the final list and
// cache some calculations.
type ItemProgression struct {
	prev      *ItemProgression // The previous ItemProgression in the path that got us here.
	Item      Item
	ValueSum  int // Total value of this progression so far
	WeightSum int // Total cost of this progression so far
}

func NewItemProgression(prev *ItemProgression, it Item) *ItemProgression {
	if it == nil {
		return nil
	}
	ret := ItemProgression{Item: it, ValueSum: it.Value(), WeightSum: it.Weight()}
	if prev != nil {
		ret.prev = prev
		ret.ValueSum += prev.ValueSum
		ret.WeightSum += prev.WeightSum
	}
	return &ret
}

func (ip *ItemProgression) Items() []Item {
	var items []Item
	if ip == nil {
		return items
	}
	if ip.prev != nil {
		items = append(items, ip.prev.Items()...)
	}
	items = append(items, ip.Item)
	return items
}

// Multiple-choice knapsack is a version of knapsack where n categories are provided, and one item must be chosen for each category.
// Each category comes with a list of items.  The algoirthm proceeds by identifying at each weight what the best item from that list
// would be. For the base-case with only one category, the most valuable item under the given weight is chosen.  To build up the
// solution, the next level must consider the best selection at a given weight to be the best one combining the current categories
// values with the best available of the previous category.
func MultiChoiceKnapsack(categories [][]Item, capacity int) []Item {
	// By cateogry, weight, and finally list of items to purchase for a given category at the given weight.
	var dynamic [][]*ItemProgression
	if len(categories) == 0 {
		return []Item{}
	}

	// Initialize our first category (category 0)
	dynamic = append(dynamic, []*ItemProgression{})
	for _, it := range initializeBaseItems(categories[0], capacity) {
		dynamic[0] = append(dynamic[0], NewItemProgression(nil, it))
	}

	// Now that we have base items, we can begin to build up the 'real items' returned by using the following relationship:
	// For any item, i, in a category, j, it is either used for a given weight, or it isn't: if it is used, then its 'value' is:
	// f(w,j) = V(i) + f(w - iw, j - 1)
	for j := 1; j < len(categories); j++ {
		// We can not have an item for 0 weight, so we initialize this as our base-case for each category.
		dynamic = append(dynamic,[]*ItemProgression{nil})
		for w := 1; w <= capacity; w++ {
			// At a minium, when increasing our weight, we can afford the item from weight - 1
			currBest := dynamic[j][w-1]
			for _, it := range categories[j] {
				// Check to see if this item could be added as part of a 'solution' for this weight. To do so it must satisfy the following
				// requirements:
				//   - Item must weigh less than the current weight limit (be affordable)
				//   - There must be a valid solution for the categories leading up to this one: d[j-1][w-it.weight]
				//   - The value using this item must be better than the current best value.
				if it.Weight() <= w {
					prevCandidate := dynamic[j-1][w-it.Weight()]
				 	if prevCandidate != nil && (currBest == nil || prevCandidate.ValueSum+it.Value() > currBest.ValueSum) {
						currBest = NewItemProgression(prevCandidate, it)
					}
				}
			}
			// If there were not previous best solution, this remains nil, if there is no improvement over the previous weight, it
			// remains the same path, otherwise we update to the newest one.
			dynamic[j] = append(dynamic[j], currBest)
		}
	}

	// We've populated our dynamic array. Return the list contained at our target capacity for the final category.
	return dynamic[len(categories)-1][capacity].Items()
}

// Return the item which gives the maximum value for a given weight.  No extra smarst here, such as finding
// the smallest or heaviest item to avoid processing.
// TODO: add such optimizations.
func initializeBaseItems(items []Item, capacity int) []Item {

	// Start with nil for the 0-weight case (no item chosen with no capacity)
	bestItems := []Item{nil}
	for i := 1; i <= capacity; i++ {
		// At a minium, when increasing our weight, we can afford the item from weight - 1
		currItem := bestItems[i-1]
		for _, it := range items {
			if it.Weight() <= i && (currItem == nil || it.Value() > currItem.Value()) {
				currItem = it
			}
		}
		bestItems = append(bestItems, currItem)
	}
	// Now we should have the best element for each price.
	return bestItems
}
