package telemetry

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service tracks telemetry events
type Service struct {
	NoTelemetry bool
	userID      string
	executionID string
	version     string
	command     string
	tracker     Tracker
}

// Setup sets up the tracker for this service
func (s *Service) Setup(command string) {
	s.command = command
	s.executionID = primitive.NewObjectID().Hex()
	if s.NoTelemetry {
		s.tracker = &noopTracker{}
	} else {
		s.tracker = newSegmentTracker()
	}
}

// SetUser sets the userID for this service to track
func (s *Service) SetUser(userID string) {
	s.userID = userID
}

// SetVersion sets the realm-cli version for this service
func (s *Service) SetVersion(version string) {
	s.version = version
}

// TrackEvent tracks the event based on the tracker
func (s *Service) TrackEvent(eventType EventType, data ...EventData) {
	if s.tracker != nil {
		s.tracker.Track(event{
			id:          primitive.NewObjectID().Hex(),
			eventType:   eventType,
			userID:      s.userID,
			time:        time.Now(),
			executionID: s.executionID,
			version:     s.version,
			command:     s.command,
			data:        data,
		})
	}
}

// Close shuts down the Service
func (s *Service) Close() {
	s.tracker.Close()
}
