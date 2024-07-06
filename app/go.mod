module app

go 1.20

require (
    github.com/aws/aws-lambda-go v1.47.0
    github.com/stretchr/testify v1.7.2
)

replace (
    golang.org/x/net => golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4
    golang.org/x/sys => golang.org/x/sys v0.0.0-20220420173948-351e6b14d86b
)
