package counters

import (
	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
)

const (
	PositionCenter = iota
	PositionTopLeft
	PositionTopCenterLeft
	PositionTopCenter
	PositionTopCenterRight
	PositionTopRight
	PositionRightTop
	PositionRight
	PositionRightBottom
	PositionBottomRight
	PositionBottomCenterRight
	PositionBottomCenter
	PositionBottomCenterLeft
	PositionBottomLeft
	PositionLeftBottom
	PositionLeft
	PositionLeftTop
)

type Positioner struct{}

func (p *Positioner) GetAnchorPointsAndMaxWidth(pos int, def *Settings) (float64, float64, float64, error) {
	counterWidth := float64(def.Width)
	offsetX := *def.Margins +
		*def.BorderWidth +
		*def.XShift +
		*def.StrokeWidth
	availableWidth := counterWidth - (offsetX * 2)

	switch pos {
	case 0:
		// Center of the counter, space usually for image
		return 0.5, 0.5, availableWidth, nil
	case 1:
		// Top left corner (10:30)
		return 0, 0, (availableWidth) / 5, nil
	case 2, 3, 4:
		// Top left, at the right of the top-left corner
		// Top, horizontally centered
		// Top right, at the left of the top-right corner
		return 0.5, 0, (availableWidth) / 5, nil
	case 5:
		// Top right corner
		return 1, 0, (availableWidth) / 5, nil
	case 6, 7, 8:
		// Right slightly up, under the top right corner, over the vertical center
		// Right, vertically centered
		// Right slightly down, over the bottom right corner, under the vertical center
		return 1, 0.5, (availableWidth) / 5, nil
	case 9:
		// Bottom right corner
		return 1, 1, (availableWidth) / 5, nil
	case 10, 11, 12:
		// Bottom right, at the left of the bottom-right corner
		// Bottom, horizontally centered
		// Bottom left, at the right of the bottom-left corner
		return 0.5, 1, (availableWidth) / 5, nil
	case 13:
		// Bottom left corner
		return 0, 1, (availableWidth) / 5, nil
	case 14, 15, 16:
		// Left down, over the bottom-left corner
		// Left, vertically centered
		// Left up, under the upper-left corner
		return 0, 0.5, (availableWidth) / 5, nil
	default:
		log.WithField("position", pos).Error("Position unknown. Valid positions are from 0 to 16 (both included)")
		return 0, 0, 0, errors.New("the position provided is not in the range 0-16")
	}
}

func (p *Positioner) GetXYPosition(pos int, def *Settings) (float64, float64, error) {
	counterWidth := float64(def.Width)
	counterHeight := float64(def.Height)

	offsetX := *def.Margins + *def.BorderWidth + *def.StrokeWidth
	offsetY := *def.Margins + *def.BorderWidth + *def.StrokeWidth

	centerVertical := counterHeight / 2
	centerHorizontal := counterWidth / 2

	topLeft := (counterWidth * 0.25) + (offsetX / 2)
	topRight := (counterWidth * 0.75) - (offsetX / 2)

	leftAbove := (counterHeight * 0.25) + (offsetY / 2)
	leftBelow := (counterHeight * 0.75) - (offsetY / 2)

	switch pos {
	case 0:
		// Center of the counter, space usually for image
		return centerHorizontal + *def.XShift, centerVertical + *def.YShift, nil
	case 1:
		// Top left corner
		return offsetX + *def.XShift, offsetY + *def.YShift, nil
	case 2:
		// Top left, at the right of the top-left corner
		return topLeft + *def.XShift, offsetY + *def.YShift, nil
	case 3:
		// Top center
		return centerHorizontal + *def.XShift, offsetY + *def.YShift, nil
	case 4:
		// Top right, at the left of the top-right corner
		return topRight + *def.XShift, offsetY + *def.YShift, nil
	case 5:
		// Top right corner
		return counterWidth - offsetX + *def.XShift, offsetY + *def.YShift, nil
	case 6:
		// Right slightly up, under the top right corner
		return counterWidth - offsetX + *def.XShift, leftAbove + *def.YShift, nil
	case 7:
		// Right
		return counterWidth - offsetX + *def.XShift, centerVertical + *def.YShift, nil
	case 8:
		// Right slightly down, over the bottom right corner
		return counterWidth - offsetX + *def.XShift, leftBelow + *def.YShift, nil
	case 9:
		// Bottom right corner
		return counterWidth - offsetX + *def.XShift, counterHeight - offsetY + *def.YShift, nil
	case 10:
		// Bottom right, at the left of the bottom-right corner
		return topRight + *def.XShift, counterHeight - offsetY + *def.YShift, nil
	case 11:
		// Bottom
		return centerHorizontal + *def.XShift, counterHeight - offsetY + *def.YShift, nil
	case 12:
		//Bottom left, at the right of the bottom-left corner
		return topLeft + *def.XShift, counterHeight - offsetY + *def.YShift, nil
	case 13:
		//Bottom left corner
		return offsetX + *def.XShift, counterHeight - offsetY + *def.YShift, nil
	case 14:
		// Left down, over the bottom-left corner
		return offsetX + *def.XShift, leftBelow + *def.YShift, nil
	case 15:
		// Left
		return offsetX + *def.XShift, centerVertical + *def.YShift, nil
	case 16:
		// Left up, under the upper-left corner
		return offsetX + *def.XShift, leftAbove + *def.YShift, nil
	default:
		log.WithField("position", pos).Error("Position unknown. Valid positions are from 0 to 16 (both included)")
		return 0, 0, errors.New("the position provided is not in the range 0-16")
	}
}

func (p *Positioner) getObjectPositions(pos int, def *Settings) (float64, float64, float64, float64, error) {
	ax, ay, _, err := p.GetAnchorPointsAndMaxWidth(pos, def)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	x, y, err := p.GetXYPosition(pos, def)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return x, y, ax, ay, nil
}
