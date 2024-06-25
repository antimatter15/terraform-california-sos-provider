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
