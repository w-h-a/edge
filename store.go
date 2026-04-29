package edge

// Store is the main type applications use.
// Open, Put, Get, Watch, Close. That is the public API.
// Store delegates to the service layer. It never touches ports directly.
type Store struct {
}
