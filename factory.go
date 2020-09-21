package main

import (
	"github.com/Gregmus2/simple-engine/common"
	"github.com/Gregmus2/simple-engine/graphics"
)

type ObjectFactory struct {
	Cfg   *common.Config
	Prog  *graphics.ProgramFactory
	Shape *graphics.ShapeHelper
}

func NewObjectFactory(cfg *common.Config, p *graphics.ProgramFactory, s *graphics.ShapeHelper) *ObjectFactory {
	return &ObjectFactory{Cfg: cfg, Prog: p, Shape: s}
}
