package service

import "context"

func (s *PlaylistService) RemoveTrack(ctx context.Context, playlistID, trackID, userID int64) error {
	p, err := s.repo.GetByID(ctx, playlistID)
	if err != nil {
		return err
	}

	if p.UserID != userID {
		return ErrForbiddenPlaylist
	}

	return s.repo.RemoveTrack(ctx, playlistID, trackID)
}
