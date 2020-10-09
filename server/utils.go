package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

// Check if the user has sysadmins rights
func isSysadmin(p *Plugin, userID string) bool {
	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		p.API.LogError("Unable to get user", "err", appErr)
		return false
	}

	return strings.Contains(user.Roles, "system_admin")
}

// Check if the user has the "delete_post" permission
func canDeletePost(p *Plugin, userID string, channelID string) bool {
	return p.API.HasPermissionTo(userID, model.PERMISSION_DELETE_POST) ||
		p.API.HasPermissionToChannel(userID, channelID, model.PERMISSION_DELETE_POST)
}

// Check if the user has the "delete_others_posts" permission
func canDeleteOthersPosts(p *Plugin, userID string, channelID string) bool {
	return p.API.HasPermissionTo(userID, model.PERMISSION_DELETE_OTHERS_POSTS) ||
		p.API.HasPermissionToChannel(userID, channelID, model.PERMISSION_DELETE_OTHERS_POSTS)
}

// Return "s" if the given number is > 1
func getPluralChar(number int) string {
	if 1 < number {
		return "s"
	}

	return ""
}

// Simplified version of SendEphemeralPost, send to the userID defined
func (p *Plugin) sendEphemeralPost(userID string, channelID string, message string) *model.Post {
	return p.API.SendEphemeralPost(
		userID,
		&model.Post{
			UserId:    p.botUserID,
			ChannelId: channelID,
			Message:   message,
		},
	)
}

// Wrapper of p.sendEphemeralPost() to one-line the return statements when a *model.CommandResponse is expected
func (p *Plugin) respondEphemeralResponse(args *model.CommandArgs, message string) *model.CommandResponse {
	_ = p.sendEphemeralPost(args.UserId, args.ChannelId, message)
	return &model.CommandResponse{}
}

// Check that a postID in form of the direct ID or a link to the post is correct
// If so, returns the postID
// If incorrect, returns "" and an error
func transformToPostID(p *Plugin, postIDToParse string, channelID string) (string, error) {
	if strings.HasPrefix(postIDToParse, "http") {
		// TODO: This is a link: transform it in a postID
		return "", errors.Errorf("Sorry, links are not supported for the moment. Please use the postID")
	}

	post, appErr := p.API.GetPost(postIDToParse)
	if appErr != nil {
		// TODO change message if internal error or user unknown
		p.API.LogError("Unable to get post", "appError :", appErr.ToJson())
		return "", errors.Errorf("unknown post `%s`", postIDToParse)
	}

	if post.ChannelId != channelID {
		return "", errors.Errorf("post `%s` is not in this channel", postIDToParse)
	}

	return post.Id, nil
}
