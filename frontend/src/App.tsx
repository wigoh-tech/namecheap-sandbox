import { Router, Routes, Route } from "react-router-dom";
import Sidebar from "./components/Sidebar";
import CheckDomain from "./pages/checkDomain";
import  BuyDomain  from "./pages/BuyDomainForm";
import Purchased from "./pages/Domainlist";

import EditDNSModal from "./components/EditDNSModal";
import BuyDomainForm from "./pages/BuyDomainForm";

export default function App() {
  return (
    
      <div className="flex min-h-screen">
        <Sidebar />
        <main className="flex-1 p-6 bg-gray-100 overflow-y-auto">
          <Routes>
            <Route path="/" element={<CheckDomain />} />
            <Route path="/buy" element={<BuyDomain domain={""} price={{
              base: 0,
              tax: 0,
              total: 0
            }} onSuccess={function (): void {
              throw new Error("Function not implemented.");
            } } />} />
            <Route path="/purchased" element={<Purchased />} />
          </Routes>
          <EditDNSModal />
        </main>
      </div>
    
  );
}
