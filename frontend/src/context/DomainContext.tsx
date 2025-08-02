import { createContext, useContext, useState } from "react";

type Price = {
  base: number;
  tax: number;
  total: number;
};

type DomainContextType = {
  selectedDomain: string;
  setSelectedDomain: (domain: string) => void;
  price: Price;
  setPrice: (price: Price) => void;
};

const DomainContext = createContext<DomainContextType | undefined>(undefined);

export function DomainProvider({ children }: { children: React.ReactNode }) {
  const [selectedDomain, setSelectedDomain] = useState("");
  const [price, setPrice] = useState<Price>({ base: 0, tax: 0, total: 0 });

  return (
    <DomainContext.Provider value={{ selectedDomain, setSelectedDomain, price, setPrice }}>
      {children}
    </DomainContext.Provider>
  );
}

export function useDomain() {
  const context = useContext(DomainContext);
  if (!context) throw new Error("useDomain must be used within DomainProvider");
  return context;
}
