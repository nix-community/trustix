// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package constants

const (
	CaValueBucket = "CaValues"

	// The root level of a given log
	LogBucket = "Log"

	// Submission queue for a log
	QueueBucket   = "Queue"
	QueueMetaBlob = "Meta"

	// The sub trees of a log
	HeadBlob      = "Head"
	VLogBucket    = "VLog"
	MapBucket     = "Map"
	VMapLogBucket = "VMapLog"
)
