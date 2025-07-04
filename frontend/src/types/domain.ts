export interface Domain {
  id: number;
  name: string;
  customer: string;
  aRecord: string;
  cName: string;
  purchased: boolean;
  revoked: boolean;
  createdAt: string;
}
export interface BuyDomainFormData {
  domain: string;
  firstName: string;
  lastName: string;
  email: string;
  address: string;
  city: string;
  phone: string;
  postalCode: string;
  country: string;
}