package test

import (
	"github.com/nirasan/go-repository-base"
	"google.golang.org/appengine/aetest"
	"testing"
)

func BenchmarkUserRepository_Find(b *testing.B) {
	ctx, f, err := aetest.NewContext()
	if err != nil {
		b.Fatal(err)
	}
	defer f()
	r := NewUserRepository(ctx)

	b.N = 1000
	for i := 0; i < b.N; i++ {
		e1 := &User{}
		r.Create(e1)

		e2, _ := r.Find(e1.ID)

		r.FindAll()

		r.Delete(e2)
	}
}

func BenchmarkDatastoreRepository_Find(b *testing.B) {
	ctx, f, err := aetest.NewContext()
	if err != nil {
		b.Fatal(err)
	}
	defer f()

	r, err := go_repository_base.NewDatastoreRepository(ctx, func() interface{} { return &User{} }, func() interface{} { return []*User{} })
	if err != nil {
		b.Fatal(err)
	}

	b.N = 1000
	for i := 0; i < b.N; i++ {
		e1 := &User{}
		r.Create(e1)

		e2, _ := r.Find(e1.ID)

		r.FindAll()

		r.Delete(e2)
	}
}
