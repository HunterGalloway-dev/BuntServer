print("START #####################################################")
db.createCollection('projects');
db.projects.insertMany([
  {
    title: "Example Project 1 Title",
    shortDescription: "Example Project 1 Short Description",
    githubURL: "Example Project 1 Github URL",
    youtubeURL: "Example Project 1 Youtube URL",
    imageURL: "Example Project 1 Image URL",
    description: "Example Project 1 Description",
    tags: ["Tag 1", "Tag 2", "Tag 3"]
  },
  {
    title: "Example Project 2 Title",
    shortDescription: "Example Project 2 Short Description",
    githubURL: "Example Project 2 Github URL",
    youtubeURL: "Example Project 2 Youtube URL",
    imageURL: "Example Project 2 Image URL",
    description: "Example Project 2 Description",
    tags: ["Tag 1", "Tag 2", "Tag 3"]
  },
  {
    title: "Example Project 3 Title",
    shortDescription: "Example Project 3 Short Description",
    githubURL: "Example Project 3 Github URL",
    youtubeURL: "Example Project 3 Youtube URL",
    imageURL: "Example Project 3 Image URL",
    description: "Example Project 3 Description",
    tags: ["Tag 1", "Tag 2", "Tag 3"]
  }
]);

print("END   #####################################################")