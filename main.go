package main

import (
	"encoding/binary"
	"io"
)

type RawImage []byte
type Label uint8

// Assumes input files are like MNIST database files
const imageSide = 28 // 28*28 pixels per image

func readImageFile(r io.Reader) (imgs []RawImage, err error) {
	var magic int32
	err = binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return
	}

	var numImages int32
	err = binary.Read(r, binary.BigEndian, &numImages)
	if err != nil {
		return
	}

	imgs = make([]RawImage, numImages)
	for i := range imgs {
		imgs[i] = make(RawImage, imageSide*imageSide)
		_, err = io.ReadFull(r, imgs[i])
		if err != nil {
			return
		}
	}
	return
}

func readLabelFile(r io.Reader) (labels []Label, err error) {
	var magic int32
	err = binary.Read(r, binary.BigEndian, &magic)
	if err != nil {
		return
	}

	var numLabels int32
	err = binary.Read(r, binary.BigEndian, &numLabels)
	if err != nil {
		return
	}

	labels = make([]Label, numLabels)
	for i := range labels {
		var label uint8
		err = binary.Read(r, binary.BigEndian, &label)
		if err != nil {
			return
		}
		labels[i] = Label(label)
	}
	return
}

func main() {
	//You would typically open your files here and pass them to the read functions
}
