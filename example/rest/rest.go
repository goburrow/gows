// Copyright 2015 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package main

import (
	"fmt"
	"os"

	"github.com/goburrow/gomelon"
	"github.com/goburrow/gomelon/core"
	"github.com/goburrow/gomelon/rest"
	"golang.org/x/net/context"
)

type User struct {
	Name string
}

// REST resource.
type resource struct {
}

func (r *resource) Path() string {
	return "/user/:name"
}

func (r *resource) GET(c context.Context) (interface{}, error) {
	params, _ := rest.PathParamsFromContext(c)
	return &User{Name: params["name"]}, nil
}

func (r *resource) POST(c context.Context) (interface{}, error) {
	return &User{}, nil
}

// Main application.
type application struct {
	rest.Application
}

func (app *application) Run(conf interface{}, env *core.Environment) error {
	if err := app.Application.Run(conf, env); err != nil {
		return err
	}
	env.Server.Register(&resource{})
	return nil
}

func main() {
	app := &application{}
	app.SetName("rest")

	err := gomelon.Run(app, os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%#v\n", err)
		os.Exit(1)
	}
}