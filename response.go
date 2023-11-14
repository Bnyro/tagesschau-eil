package main

import "time"

type News struct {
	SophoraID   string    `json:"sophoraId"`
	ExternalID  string    `json:"externalId"`
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	TeaserImage struct {
		Alttext       string `json:"alttext"`
		ImageVariants struct {
			Square144 string `json:"1x1-144"`
			Square256 string `json:"1x1-256"`
			Square432 string `json:"1x1-432"`
			Square640 string `json:"1x1-640"`
			Square840 string `json:"1x1-840"`
			Land256   string `json:"16x9-256"`
			Land384   string `json:"16x9-384"`
			Land512   string `json:"16x9-512"`
			Land640   string `json:"16x9-640"`
			Land960   string `json:"16x9-960"`
			Land1280  string `json:"16x9-1280"`
			Land1920  string `json:"16x9-1920"`
		} `json:"imageVariants"`
		Type string `json:"type"`
	} `json:"teaserImage,omitempty"`
	Tags []struct {
		Tag string `json:"tag"`
	} `json:"tags"`
	UpdateCheckURL string `json:"updateCheckUrl"`
	Tracking       []struct {
		Sid  string `json:"sid"`
		Src  string `json:"src"`
		Ctp  string `json:"ctp"`
		Pdt  string `json:"pdt"`
		Otp  string `json:"otp"`
		Cid  string `json:"cid"`
		Pti  string `json:"pti"`
		Bcr  string `json:"bcr"`
		Type string `json:"type"`
	} `json:"tracking"`
	Topline       string        `json:"topline,omitempty"`
	FirstSentence string        `json:"firstSentence,omitempty"`
	Details       string        `json:"details,omitempty"`
	Detailsweb    string        `json:"detailsweb,omitempty"`
	ShareURL      string        `json:"shareURL,omitempty"`
	Geotags       []interface{} `json:"geotags,omitempty"`
	RegionID      int           `json:"regionId,omitempty"`
	RegionIds     []interface{} `json:"regionIds,omitempty"`
	Type          string        `json:"type"`
	BreakingNews  bool          `json:"breakingNews,omitempty"`
	Streams       struct {
		H264S             string `json:"h264s"`
		H264M             string `json:"h264m"`
		H264Xl            string `json:"h264xl"`
		Adaptivestreaming string `json:"adaptivestreaming"`
	} `json:"streams,omitempty"`
	Alttext       string `json:"alttext,omitempty"`
	Copyright     string `json:"copyright,omitempty"`
	BrandingImage struct {
		Title         string `json:"title"`
		Copyright     string `json:"copyright"`
		Alttext       string `json:"alttext"`
		ImageVariants struct {
			Original string `json:"original"`
		} `json:"imageVariants"`
		Type string `json:"type"`
	} `json:"brandingImage,omitempty"`
	Ressort  string `json:"ressort,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type NewsResponse struct {
	News                []News        `json:"news"`
	Regional            []interface{} `json:"regional"`
	NewStoriesCountLink string        `json:"newStoriesCountLink"`
	Type                string        `json:"type"`
	NextPage            string        `json:"nextPage"`
}
