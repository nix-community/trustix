package main

// import (
// 	"fmt"
// )

// // Implement MapStore interface from smt library
// // type MapStore interface {
// // 	Get(key []byte) ([]byte, error)
// // 	Set(key []byte, value []byte) error
// // 	Delete(key []byte) error
// // }

// // InvalidKeyError is thrown when a key that does not exist is being accessed.
// type InvalidKeyError struct {
// 	Key []byte
// }

// func (e *InvalidKeyError) Error() string {
// 	return fmt.Sprintf("invalid key: %s", e.Key)
// }

// type SimpleMap struct {
// 	m map[string][]byte
// }

// func NewSimpleMap() *SimpleMap {
// 	return &SimpleMap{
// 		m: make(map[string][]byte),
// 	}
// }

// func (sm *SimpleMap) Get(key []byte) ([]byte, error) {
// 	if value, ok := sm.m[string(key)]; ok {
// 		return value, nil
// 	}
// 	return nil, &InvalidKeyError{Key: key}
// }

// func (sm *SimpleMap) Set(key []byte, value []byte) error {
// 	sm.m[string(key)] = value
// 	return nil
// }

// func (sm *SimpleMap) Delete(key []byte) error {
// 	_, ok := sm.m[string(key)]
// 	if ok {
// 		delete(sm.m, string(key))
// 		return nil
// 	}
// 	return &InvalidKeyError{Key: key}
// }
