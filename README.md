# terraform-california-sos-provider

What if we could define our corporate infrastructure with code the way we can with cloud
infrastructure?

What if the state of our company was specified declaratively— just edit your corporate
configuration, run `terraform plan` to review the documents, and `terraform apply` to file them.

Imagine a world where you could just fork an existing corporate structure and start from there.

## Implementation

The current implementation supports filing forms LLC-1, LLC-4/8, AMDT-STK-NA with the California 
Secretary of State, for incorporating an LLC, changing the name of the company, and dissolving the 
company. It uses Lob as an API to send physical mail for the appropriate documents. 

This code was written in early 2020. To my knowledge, it has never been used end-to-end. Pull requests accepted.

## Building

```
go build -o terraform-provider-terracorp && terraform init
```

## Thoughts

Terraform has a lot of infrastructure built out which overlaps with what we'd probably need in such
a system, but after working a bit with the Terraform provider system, it's becoming clear to me that
an actual implementation of the system probably ought to be slightly more custom.

Terraform doesn't really spend much time planning things out. This makes sense for simpler CRUD
resources where they are pretty typically created or removed with not very much consequence. For
corporations, it makes sense to have a bit of a more thorough planning phase. That way we can check
invariants (for instance, when registering a corporate entity, we should check that the name does
not interfere with an existing one— we probably don't want to wait until we try to `terraform apply`
before we notice that this is a problem). There's a chance we can use
[CustomizeDiff](https://www.terraform.io/docs/extend/resources/customizing-differences.html) for
this.

We'd also like to be able to emit some more precise information about what exactly needs to happen.
