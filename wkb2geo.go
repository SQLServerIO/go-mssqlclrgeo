package mssqlclrgeo

import (
	"bytes"

	"fmt"
)

// http://edndoc.esri.com/arcsde/9.1/general_topics/wkb_representation.htm
// https://github.com/postgis/postgis/blob/2.1.0/doc/ZMSgeoms.txt
// http://www.opengeospatial.org/standards/sfs

type wkbByteOrder uint8

const (
	wkbXDR wkbByteOrder = 0 // Big Endian
	wkbNDR wkbByteOrder = 1 // Little Endian
)

const wkbZ uint32 = 0x80000000
const wkbM uint32 = 0x40000000
const wkbSRID uint32 = 0x20000000

type wkbGeometryType uint32

const (
	typeWkbPoint              wkbGeometryType = 1
	typeWkbLineString         wkbGeometryType = 2
	typeWkbPolygon            wkbGeometryType = 3
	typeWkbMultiPoint         wkbGeometryType = 4
	typeWkbMultiLineString    wkbGeometryType = 5
	typeWkbMultiPolygon       wkbGeometryType = 6
	typeWkbGeometryCollection wkbGeometryType = 7
	typeWkbCircularString     wkbGeometryType = 8
	typeWkbCompoundCurve      wkbGeometryType = 9
	typeWkbCurvePolygon       wkbGeometryType = 10
	typeWkbMultiCurve         wkbGeometryType = 11
)

type linearRing struct {
	numPoints uint32
	points    []Point
}

type wkbPoint struct {
	wkbGeometry WkbGeometry
	point       Point
}

type wkbLineString struct {
	wkbGeometry WkbGeometry
	numPoints   uint32
	points      []Point
}
type wkbPolygon struct {
	wkbGeometry WkbGeometry
	numRings    uint32
	rings       []linearRing
}
type wkbMultiPoint struct {
	wkbGeometry   WkbGeometry
	num_wkbPoints uint32
	wkbPoints     []wkbPoint
}
type wkbMultiLineString struct {
	wkbGeometry        WkbGeometry
	num_wkbLineStrings uint32
	wkbLineStrings     []wkbLineString
}
type wkbMultiPolygon struct {
	wkbGeometry     WkbGeometry
	num_wkbPolygons uint32
	wkbPolygons     []wkbPolygon
}
type wkbGeometryCollection struct {
	wkbGeometry       WkbGeometry
	num_wkbGeometries uint32
	wkbGeometries     []interface{}
}

type WkbGeometry struct {
	byteOrder wkbByteOrder
	wkbType   wkbGeometryType
	srid      uint32
	//geo       interface{}

	hasZ    bool
	hasM    bool
	hasSRID bool
}

func ParseWkb(data []byte) (g interface{}, err error) {

	var buffer = bytes.NewBuffer(data[0:])
	geom, err := readWKBGeometry(buffer)

	return geom, err

}

func WkbToUdtGeo(geom interface{}) (data []byte, err error) {

	var b Builder

	switch wkbGeom := geom.(type) {

	case *wkbPoint:
		b.Srid = wkbGeom.wkbGeometry.srid
		b.AddShape(SHAPE_POINT)
		b.AddFeature()
		b.AddPoint(wkbGeom.point.X, wkbGeom.point.Y, wkbGeom.point.Z, wkbGeom.point.M)

		fmt.Println("WkbPoint 1")

	case *wkbLineString:
		fmt.Println("WkbLineString 2")
		b.Srid = wkbGeom.wkbGeometry.srid
		b.AddShape(SHAPE_LINESTRING)
		b.AddFeature()
		for _, point := range wkbGeom.points {
			b.AddPoint(point.X, point.Y, point.Z, point.M)
		}
	case *wkbPolygon:
		b.Srid = wkbGeom.wkbGeometry.srid
		b.AddShape(SHAPE_POLYGON)
		for _, ring := range wkbGeom.rings {
			b.AddFeature()
			for _, point := range ring.points {
				b.AddPoint(point.X, point.Y, point.Z, point.M)
			}
		}
		fmt.Println("WkbPolygon 3")
	case *wkbMultiPoint:
		//fmt.Println("MultiPoint 4")
	case *wkbMultiLineString:
		//fmt.Println("MultiLineString 5")
	case *wkbMultiPolygon:
		fmt.Println("MultiPolygon 6")
		fmt.Printf("num_wkbPolygons: %d", wkbGeom.num_wkbPolygons)

	case *wkbGeometryCollection:
		//fmt.Print("GeometryCollection 7")
	default:
		fmt.Print("other")

	}

	return b.Generate()
}
