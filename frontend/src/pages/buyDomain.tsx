import BuyDomainForm from "../pages/BuyDomainForm";

export default function BuyDomainPage() {
  return <BuyDomainForm onSuccess={() => alert("Success!")} />;
}
