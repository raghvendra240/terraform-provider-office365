# Terraform provider for office365

This provider allows to perform Create ,Read ,Update ,Delete ,activate , deactivate and Import user of office365

## Requirements

terraform office365 provider is based on Terraform, this means that you need

1. [Go](https://golang.org/doc/install) 1.11 or later
2. [Terraform](https://www.terraform.io/downloads.html) v0.13.0 or later
3. [office365 Developer account](https://developer.microsoft.com/en-us/microsoft-365/dev-program)


## Setup office365 account

1. Create office365 developer account by following [this](https://docs.microsoft.com/en-us/office/developer-program/microsoft-365-developer-program)
2. Set up a Microsoft 365 developer sandbox subscription by following [this](https://docs.microsoft.com/en-us/office/developer-program/microsoft-365-developer-program-get-started)
3. Register your application in Azure AD by following [this](https://docs.microsoft.com/en-us/office/office-365-management-api/get-started-with-office-365-management-apis#register-your-application-in-azure-ad)
4. Generate a new key for your application by following [this](https://docs.microsoft.com/en-us/office/office-365-management-api/get-started-with-office-365-management-apis#generate-a-new-key-for-your-application)


## Initialise Expensify Provider in local machine 
1. Clone the repository  to $GOPATH/src/github.com/office365/terraform-provider-office365 <br>
2. Add the partnerUserID and partnerUserSecret generated to respective fields in `main.tf` <br>
3. Run the following command :
```golang
go mod init terraform-provider-expensify
go mod tidy
````

```


## Installation

1. Clone this repository
2. Run the following command to create a vendor subdirectory which will comprise of all provider dependencies.
     ```mkdir -p %APPDATA%//.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/[architecture name]/```
     you can find full list of architecture names [here](https://golang.org/doc/install/source#environment)
    for example:
    ```
    mkdir -p %APPDATA%/terraform.d/plugins/raghvnedra/clvertap/office365/0.1/windows_amd64
    ```
2. Run ``` go build -o terraform-provider-office365.exe ``` This will save the binary (.exe) file in the main/root directory.
3. Run this command to move this binary file to appropriate location.
```
move terraform-provider-office365.exe %APPDATA%\terraform.d\plugins\raghvendra\clevertap\office365\0.1\[OS_ARCH]
```
5. Otherwise you can manually move the file from current directory to destination directory.



## Comands to Run the Provider

1. terraform init
2. terraform plan
3. terrafrom apply (To create or update the user)
4. terraform destroy (To destroy the created user)

## Steps to run import command

1. Write manually a resource configuration block for the resource, to which the imported object will be mapped.

2. RUN `` terraform import office365_user_manage.<block_name> <userPrincipalName> ``

3. Check for the attributes in the .tfstate file and fill them accordingly in resource block.

### Usage Example
```
# This is required for Terraform 0.13+
terraform {
  required_providers {
    office365 = {
      version = "~> 0.1"
      source  = "Raghvendra/clevertap/office365"
    }
  }
}

# configure provider
provider "office365" {
   client_id     = "---"
   client_secret = "---"
   tenant_id     ="---"
  
  #Above parameters are required to access the developer account
}

#Creating User
resource "office365_user_manage" "example" {
   display_name        ="---"
   user_principal_name ="example@<officce365domain>.onmicrosoft.com"
   password            ="---"
   account_enabled     ="---"
}

#Updaing User
 - To update any parameter just  change them in the respective resouce block and run terraform apply
 - passowrd and user_principal_name can not be changed.
 - To activate or deactivate user, set true/false to account_enabled accordingly.


#Get user Information
data "office365_users" "example" {
    userprincipalname ="example@<officce365domain>.onmicrosoft.com"
}
 output "example" {
      value          = data.office365_users.example
}

```

### Supported Arguments in resource block
- ``display_name`` (required) The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name.<br/>

- ``user_principal_name`` (required)The user principal name (UPN) of the user. The UPN is an Internet-style login name for the user based on the Internet standard RFC 822. By convention, this should map to the user's email name. The general format is alias@domain, where domain must be present in the tenant's collection of verified domains. This property is required when a user is created. 

- ``password`` (required) password for created user.<br/>

- ``account_enable`` (optional) Takes only true/false as input and allow us to activate/deactivate account(Default:true).

- ``given_name`` (optional) first name of user.

- ``surname`` (optional) second name of user.

- ``Job-title`` (optional) The user's job title. Maximum length is 128 characters.

- ``office_location``(optional) office location of user.

- ``mail_nick_name`` (optional) The mail alias for the user.

- ``mobile_phone`` (optional) The primary cellular telephone number for the user.
- ``postal_code`` (optional)The postal code for the user's postal address. The postal code is specific to the user's country/region.

- ``state`` (optional)The state or province in the user's address.
- ``preferred_language`` (optional)The preferred language for the user. Should follow ISO 639-1 Code; for example "en-US".
- ``street_address`` (optional)The street address of the user's place of business.
- ``usage_location`` (optional)A two letter country code (ISO standard 3166). Required for users that will be assigned licenses due to legal requirement to check for availability of services in countries. Examples include: "US", "JP", and "GB".

