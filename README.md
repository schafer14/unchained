# Problem

God's Unchained applications are becoming more and more complex as we add more stuff. 
In addition we have feature drift between repositories based on when they were written. 

# Solution Parameters

We need any solutions to be:

- Simple for developers
- Controllable and configurable so all applications are kept in line

# Proposed solution 

## TL DR 

Use buildpacks to control the runtime environment of our applications.

## More 

We can radically simplify our applications while maintaining the highest levels of 
control by creating a simple runtime that is specific to our needs (mostly handling
http and queue events).

## Benefits

- Control
- Simplicity
- Reduced boilerplate

# Example Implementation

An example implementation on what a handler would look like is located in the `/sample`
directory of this repository. This example handles HTTP requests but the same idea
could be applied to different types of handlers.

The `fns` directory contains a list of handlers. These have been reduced to the simplest
state (Go functions). There is no need to worry about transport logic. 

Dependencies can be injected by creating provider functions in the `dependencies` directory.
Any dependency is available to a handler at runtime. 

## Running the example

### Dependencies

- Docker
- [pack](https://buildpacks.io/docs/tools/pack/)

### Running

```console
# This line may be unnecessary ?!
pack config default-builder cnbs/sample-builder:bionic

pack build test-unchained --path ./sample --clear-cache --buildpack ./unchained
```

This will create a docker image. To run the image:

```console
docker run -it -p 8080:8080 test-unchained
```

And now there is a server bound to port 8080.


```console
curl -i localhost:8080/assets/new -d '{"proto": 3, "quality": 3, "walletAddress": "132asdf"}'
```

# Thoughts

There is a lot happening here that warrants discussion, but before we talk about how
this is actually implemented it would be good to consider whether we even want this
sort of thing. 

Also don't look into the `unchained` directory. It is a mess, but could be cleaned up
pretty easily. 

