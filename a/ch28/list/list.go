// Package list provides generic implementations of singly and doubly linked lists.
package list

// --- Doubly Linked List ---

// DoublyNode is a node in a doubly linked list.
type DoublyNode[T any] struct {
	Value T
	prev  *DoublyNode[T]
	next  *DoublyNode[T]
}

// Prev returns the previous list node or nil.
func (n *DoublyNode[T]) Prev() *DoublyNode[T] { return n.prev }

// Next returns the next list node or nil.
func (n *DoublyNode[T]) Next() *DoublyNode[T] { return n.next }

// DoublyLinkedList is a generic doubly linked list.
type DoublyLinkedList[T any] struct {
	head *DoublyNode[T]
	tail *DoublyNode[T]
	Len  int
}

// NewDoubly creates a new, empty doubly linked list.
func NewDoubly[T any]() *DoublyLinkedList[T] {
	return &DoublyLinkedList[T]{}
}

// Front returns the first node of the list or nil if the list is empty.
func (l *DoublyLinkedList[T]) Front() *DoublyNode[T] {
	return l.head
}

// Back returns the last node of the list or nil if the list is empty.
func (l *DoublyLinkedList[T]) Back() *DoublyNode[T] {
	return l.tail
}

// PushFront adds a new node with the given value to the front of the list.
func (l *DoublyLinkedList[T]) PushFront(value T) *DoublyNode[T] {
	node := &DoublyNode[T]{Value: value}
	l.PushFrontNode(node) // Use the helper to add the new node
	return node
}

// PushFrontNode adds an existing node to the front of the list.
// The node must be detached from any other list before calling this.
func (l *DoublyLinkedList[T]) PushFrontNode(node *DoublyNode[T]) {
	node.prev = nil // Ensure the node is detached from its previous context
	node.next = nil // Ensure the node is detached from its previous context

	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		node.next = l.head
		l.head.prev = node
		l.head = node
	}
	l.Len++
}

// Remove removes a node from the list.
func (l *DoublyLinkedList[T]) Remove(node *DoublyNode[T]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
	node.prev = nil // Detach node from the list
	node.next = nil // Detach node from the list
	l.Len--
}

// MoveToFront moves a node to the front of the list.
func (l *DoublyLinkedList[T]) MoveToFront(node *DoublyNode[T]) {
	if l.head == node {
		return // Already at the front
	}
	l.Remove(node)
	l.PushFrontNode(node) // Use the existing node
}

// --- Singly Linked List ---

// SinglyNode is a node in a singly linked list.
type SinglyNode[T any] struct {
	Value T
	next  *SinglyNode[T]
}

// Next returns the next list node or nil.
func (n *SinglyNode[T]) Next() *SinglyNode[T] { return n.next }

// SinglyLinkedList is a generic singly linked list.
type SinglyLinkedList[T any] struct {
	head *SinglyNode[T]
	tail *SinglyNode[T]
	Len  int
}

// NewSingly creates a new, empty singly linked list.
func NewSingly[T any]() *SinglyLinkedList[T] {
	return &SinglyLinkedList[T]{}
}

// Front returns the first node of the list or nil.
func (l *SinglyLinkedList[T]) Front() *SinglyNode[T] {
	return l.head
}

// PushBack adds a new node with the given value to the back of the list.
func (l *SinglyLinkedList[T]) PushBack(value T) *SinglyNode[T] {
	node := &SinglyNode[T]{Value: value}
	if l.head == nil {
		l.head = node
		l.tail = node
	} else {
		l.tail.next = node
		l.tail = node
	}
	l.Len++
	return node
}

// RemoveFront removes the first node from the list.
func (l *SinglyLinkedList[T]) RemoveFront() {
	if l.head == nil {
		return
	}
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	}
	l.Len--
}
