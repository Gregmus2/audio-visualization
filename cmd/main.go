package main

import (
	engine "github.com/Gregmus2/simple-engine"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"music/internal"
	"runtime"
)

func main() {
	runtime.LockOSThread()

	c, err := engine.BuildContainer()
	if err != nil {
		logrus.WithError(err).Fatal("error building DI container")
	}

	err = buildContainer(c)
	if err != nil {
		logrus.WithError(err).Fatal("error building DI container")
	}

	if err := c.Invoke(func(app *engine.App, music *internal.Music) {
		app.InitWithScene(music)
		app.Loop()
	}); err != nil {
		logrus.Fatal(err)
	}
}

func buildContainer(c *dig.Container) error {
	if err := c.Provide(internal.NewObjectFactory); err != nil {
		return err
	}

	if err := c.Provide(internal.NewMusic); err != nil {
		return err
	}

	return nil
}
