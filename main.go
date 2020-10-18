package main

import (
    "context"
    "encoding/json"
    "fmt"

    "demo/db"
)

func main() {
    if err := run(); err != nil {
        panic(err)
    }
}

func run() error {
    client := db.NewClient()
    err := client.Connect()
    if err != nil {
        return err
    }

    defer func() {
        err := client.Disconnect()
        if err != nil {
            panic(err)
        }
    }()

    ctx := context.Background()

    // create a post
    createdPost, err := client.Post.CreateOne(
        db.Post.Title.Set("Hi from Prisma!"),
        db.Post.Published.Set(true),
        db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
    ).Exec(ctx)
    if err != nil {
        return err
    }

    result, _ := json.MarshalIndent(createdPost, "", "  ")
    fmt.Printf("created post: %s\n", result)

    // find a single post
    post, err := client.Post.FindOne(
        db.Post.ID.Equals(createdPost.ID),
    ).Exec(ctx)
    if err != nil {
        return err
    }

    result, _ = json.MarshalIndent(post, "", "  ")
    fmt.Printf("post: %s\n", result)

    // for optional/nullable values, you need to check the function and create two return values
    // `name` is a string, and `ok` is a bool whether the record is null or not. If it's null,
    // `ok` is false, and `name` will default to Go's default values; in this case an empty string (""). Otherwise,
    // `ok` is true and `desc` will be "my description".
    name, ok := post.Desc()
    if !ok {
        return fmt.Errorf("post's name is null")
    }

    fmt.Printf("The posts's name is: %s\n", name)

    return nil
}