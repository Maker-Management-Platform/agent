import { IconBook } from "@tabler/icons-react";
import { LibMain } from "./pages/libMain/LibMain";
import { LibSettings } from "./pages/libSettings/LibSettings";
import { NavItem } from "@/navigation";
import { Asset } from "./components/asset/Asset";

export const navItem: NavItem = {
    label: 'Libraries',
    icon: IconBook,
    path: '/lib',
    element: <LibMain />,
    children: [
        {
            label: 'Libraries',
            path: ':id?',
            element: <Asset />,
        },
        {
            label: 'Library',
            path: '/lib/:id',
            element: <LibSettings />,
        },
    ]
}
