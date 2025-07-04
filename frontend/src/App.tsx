import { useState } from "react";
import DomainSearch from "./components/DomainSearch";
import DomainDetails from "./components/DomainDetails";
import BuyDomainForm from "./components/BuyDomainForm";

export default function App() {
  const [domain, setDomain] = useState("");
  const [step, setStep] = useState(1); // 1: Search, 2: Details, 3: Form

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200 px-4">
      <div className="w-full max-w-md bg-white rounded-2xl shadow-lg p-6 text-center">
        <h1 className="text-2xl font-bold text-blue-700 mb-2">Sandbox Domain Reseller</h1>
        <p className="text-gray-600 text-sm mb-6">Search, View, and Purchase your domain</p>
        

        {step === 1 && (
          <DomainSearch
            onAvailable={(domain) => {
              setDomain(domain);
              setStep(2);
            }}
          />
        )}

        {step === 2 && (
          <DomainDetails
            domain={domain}
            onBuyClick={() => setStep(3)}
          />
        )}

        {step === 3 && <BuyDomainForm domain={domain} />}
      </div>
    </div>
  );
}
