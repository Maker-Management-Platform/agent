import { IconGauge } from "@tabler/icons-react";
import { Dashboard } from "./pages/dashboard/Dashboard";
import { NavItem } from "@/navigation";

export const navItem: NavItem = {
    label: 'Dashboard',
    icon: IconGauge,
    path: '/dashboard',
    element: <Dashboard />,
}