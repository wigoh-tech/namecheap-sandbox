import { Routes, Route, Navigate } from "react-router-dom";
import Sidebar from "./components/Sidebar";
import CheckDomain from "./pages/checkDomain";
import BuyDomainForm from "./pages/BuyDomainForm"; // 👈 import the updated form (not passing props anymore)
import Purchased from "./pages/Domainlist";
import EditDNSModal from "./components/DNSModal";

export default function App() {
  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="flex-1 p-6 bg-gray-100 overflow-y-auto">
        <Routes>
           <Route path="/" element={<Navigate to="/check" />} />
          <Route path="/check" element={<CheckDomain />} />
          <Route path="/buy" element={<BuyDomainForm onSuccess={() => alert("Domain purchased")} />} /> {/* ✅ Fix */}
          <Route path="/domains" element={<Purchased />} />
        </Routes>
        <EditDNSModal />
      </main>
    </div>
  );
}
