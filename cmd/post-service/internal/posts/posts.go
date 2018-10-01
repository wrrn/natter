package posts

import (
	"github.com/google/uuid"
	pb "github.com/wrrn/natter/pkg/post"
)

var funnyPosts = []pb.Post{
	{
		Msg: "I have a lot of growing up to do. I realized that the other day inside my fort.",
	},
	{
		Msg: "The easiest time to add insult to injury is when you’re signing somebody’s cast.",
	},
	{
		Msg: "You don’t need a parachute to go skydiving. You need a parachute to go skydiving twice.",
	},
	{
		Msg: "Letting go of a loved one can be hard. But sometimes, it’s the only way to survive a rock climbing catastrophe.",
	},
	{
		Msg: "Build a man a fire, and he’ll be warm for a day. Set a man on fire, and he’ll be warm for the rest of his life.",
	},
	{
		Msg: "To steal ideas from one person is plagiarism. To steal from many is research.",
	},
	{
		Msg: "My new girlfriend works at the zoo. I think she’s a keeper.",
	},
	{
		Msg: "The inventor of the throat lozenge has died. There will be no coffin at his funeral.",
	},
	{
		Msg: "Two peanuts were walking down the street. One was a salted.",
	},
	{
		Msg: "I used to have a job at a calendar factory but I got the sack because I took a couple of days off.",
	},
	{
		Msg: "Two guys walk into a bar, the third one ducks.",
	},
	{
		Msg: "I had a dream that I was a muffler last night. I woke up exhausted!",
	},
	{
		Msg: "What is Beethoven's favorite fruit? A ba-na-na-na.",
	},
	{
		Msg: "What's Forrest Gump's password? 1forrest1",
	},
	{
		Msg: "I asked my dad for his best dad joke and he said 'You'",
	},
	{
		Msg: "I gave all my dead batteries away today… Free of charge.",
	},
	{
		Msg: "I needed a password eight characters long so I picked Snow White and the Seven Dwarfs.",
	},
	{
		Msg: "What’s the advantage of living in Switzerland? Well, the flag is a big plus.",
	},
	{
		Msg: "Just watched a documentary about beavers… It was the best damn program I’ve ever seen.",
	},
	{
		Msg: "I’m reading a book on the history of glue – can’t put it down.",
	},
	{
		Msg: "I went to the zoo the other day, there was only one dog in it. It was a shitzu.",
	},
	{
		Msg: "What did the 0 say to the 8? Nice belt.",
	},
	{
		Msg: "Why do scuba divers fall backwards into the water? Because if they fell forwards they’d still be in the boat.",
	},
	{
		Msg: "Have you ever heard of a music group called Cellophane? They mostly wrap.",
	},
	{
		Msg: "I’m a big fan of whiteboards. I find them quite re-markable.",
	},
}

func dadJokes() map[string]pb.Post {
	posts := make(map[string]pb.Post)
	for _, p := range funnyPosts {
		newPost := pb.Post{
			Uuid: uuid.New().String(),
			Msg:  p.Msg,
		}
		posts[newPost.Uuid] = newPost
	}

	return posts
}
