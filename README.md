## Overview

Your mission, should you choose to accept it, is to construct a system which uses a load balancer
to balance traffic between any number of nginx containers, which serve a static file containing the
IP addresses of all of the currently connected nginx containers.
This should all work on the fly as nginx containers are added or removed.

## Requirements

To make things more interesting, you will have some specific constraints and requirements:

1. Each nginx instance must run in a separate docker container.

2. You must use the alpine:3.10 image from the public docker repository as the base of your
docker images (yes, we know there are pre‐built nginx images).

3. This should be a highly available solution. Incorporate load balancing and container
lifecycle management tools with your choice of technologies. This may include Terraform,
Ansible, Consul, HAProxy, AWS ECS, AWS ELB, Kubernetes, Docker Compose or
similar.

4. Commit any work to a new repo, push, and share the git URL with us when you are ready
to start showing your work.

5. Make sure that there are sufficient and easy‐to‐find instructions and documentation
contained within the repository so that we can build, run, and test what you've done as
appropriate.

Bonus:
1. How would this be configured to maximize availability. (documentation only)
2. What loads would this spinup be able to handle. (documentation only)
3. How would logging, security be applied. (documentation only)

As long as your solution follows the above requirements, you are free (and encouraged!) to add
any other bits and pieces you'd like to make it fancier or otherwise more fun.

What we’ll be looking for

This is intentionally a very open ended exercise, so there is no single right or wrong answer.
We will be looking for things like:
1. Everything runs.
2. How you approached and thought through the problem.
3. How well you know and used the features of the tools and libraries you made use of.
4. How professionally "put together" and engineered everything is.

Remember though, this is a chance for you to be creative and show off, so let us see what you
can do, and try to have some fun with it!


## Discussion

### Node information

The current implementation using CronJobs will have the following issues:
- It runs only over a period of time. If real time updates are necessary, then a 
deployment might be a better option with a polling logic which is always pushing the current state. 
The downside of this approach is always hitting the blob storage endpoint. In that case, using 
a distributed storage to place a file which is mounted to all nodes would be a better approach; however,
this would be limited by cluster or region boundaries.
- Another point to take in account is that the current implementation does not keep previous data; hence, 
information is for the specific moment in time. This can be fixed if files are uploaded with a unique 
string; then, another logic in the frontend can pull the right information checking the latest timestamp.
Another approach is pulling data from the bucket, and add new data until reaches a threshold then partition 
the data.
