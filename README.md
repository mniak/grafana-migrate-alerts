Grafana Alert Rule Migration tool
==================================

Generate Terraform definitions for Alert Rules without using the Grafana API.

Useful for usage in scenarios where the user can't use `terraform import` either because they don't have adequate permissions for using Grafana API or because they don't have permissions to access the Terraform state file.


### Flags

#### `--base-url URL` (Required)
Indicates the URL of the grafana server.

Example: `https://observability.mygrafana.local`

#### `--session TOKEN` (Required)
Grafana session token. It can be obtained from the cookie `grafana_session` of an active session.

#### `--group-suffix WORD` and `--rule-suffix WORD` (Optional)
The generator will append those suffixes in the Alert Rule Group Name and on the Alert Rule Title respectively.

#### `--all` (Optional)
By default, the generator will only generate Terraform definitions for Alert Rules that were created using the Grafana web interface.

This flag removes this restriction an all Alert Rules will have Terraform Definitions generated.