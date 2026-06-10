package service

import "context"

func (s *PlaylistService) Delete(ctx context.Context, id, userID int64) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if p.UserID != userID {
		return ErrForbiddenPlaylist
	}

	return s.repo.Delete(ctx, id)
}
