package mssqlclrgeo

//https://msdn.microsoft.com/en-us/library/ee320529(v=sql.105).aspx
//https://msdn.microsoft.com/en-us/library/microsoft.sqlserver.types.sqlgeometrybuilder.begingeometry.aspx
type Builder struct {
	g          Geometry
	Srid       uint32
	stackShape *Shape
}

func (b *Builder) AddShape(shape_type SHAPE) {
	shape := &Shape{OpenGisType: shape_type}
	shape.FigureOffset = int32(len(b.g.Figures))
	if b.stackShape == nil {
		shape.ParentOffset = -1
	}
	b.g.Shapes = append(b.g.Shapes, *shape)

}
func (b *Builder) AddFeature() {
	figure := &Figure{Attribute: FIGURE_V2_LINE}

	figure.Offset = uint32(len(b.g.Points))
	b.g.Figures = append(b.g.Figures, *figure)
}
func (b *Builder) AddPoint(x float64, y float64, z float64, m float64) {
	point := &Point{X: x, Y: y, Z: z, M: m}
	b.g.Points = append(b.g.Points, *point)
}

func (b *Builder) Generate() (data []byte, err error) {
	b.g.SRID = int32(b.Srid)
	b.g.Version = 1
	return WriteGeometry(b.g, false)
}
