# terraform-california-sos-provider

What if we could define our corporate infrastructure with code the way we can with cloud
infrastructure?

What if the state of our company was specified declaratively— just edit your corporate
configuration, run `terraform plan` to review the documents, and `terraform apply` to file them.

Imagine a world where you could just fork an existing corporate structure and start from there.

The current implementation supports filing forms LLC-1, LLC-4/8, AMDT-STK-NA with the California 
Secretary of State, for incorporating an LLC, changing the name of the company, and dissolving the 
company. It uses Lob as an API to send physical mail for the appropriate documents. 

This code was written in early 2020. To my knowledge, it has never been used end-to-end. Pull requests accepted.

## Usage

To create a company, create `main.tf`

```
resource "terracorp_llc" "acme-llc" {
  name = "ACME LLC"
  owner_name = "Kevin Kwok"
  phone_number = "(123) 456 - 7891"
  address = "1 Transfinite Loop"
  city = "San Francisco"
  zip = "94107"
}
```

To see that it will file `Form LLC-1` to incorporate, run

```
terraform plan
```

To send the paperwork to the California Secretary of State at Sacramento P.O. Box 944260 over physical mail through the Lob.com API, run

```
terraform apply
```

After the company has been incorporated, you can modify the `name` field in `main.tf` and run `terraform apply` and it will file `Form AMDT-STK-NA` to rename the company. 

To dissolve the corporation you can run `terraform apply -destroy` and it will file `Form LLC-4/8`

## Building

```
go build -o terraform-provider-terracorp && terraform init
```




