package transfers

import "github.com/aws/aws-sdk-go/service/route53domains"

type Transfer struct {
	TransferInput          *route53domains.TransferDomainToAnotherAwsAccountInput
	TransferResponse       *route53domains.TransferDomainToAnotherAwsAccountOutput
	AcceptanceInput        *route53domains.AcceptDomainTransferFromAnotherAwsAccountInput
	AcceptanceResponse     *route53domains.AcceptDomainTransferFromAnotherAwsAccountOutput
	SourceOperationsInput  *route53domains.ListOperationsInput
	SourceOperationsOutput *route53domains.ListOperationsOutput
	TargetOperationsInput  *route53domains.ListOperationsInput
	TargetOperationsOutput *route53domains.ListOperationsOutput
}

func (t *Transfer) GenerateAcceptance() {
	t.SourceOperationsInput = &route53domains.ListOperationsInput{}
	t.TargetOperationsInput = &route53domains.ListOperationsInput{}
	t.AcceptanceInput = &route53domains.AcceptDomainTransferFromAnotherAwsAccountInput{
		DomainName: t.TransferInput.DomainName,
		Password:   t.TransferResponse.Password,
	}
}
