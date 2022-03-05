import dynamic from "next/dynamic";

const AdminConfigNoSSR = dynamic(() => import("../../components/admin/index"), {
  ssr: false,
});

export default AdminConfigNoSSR;
