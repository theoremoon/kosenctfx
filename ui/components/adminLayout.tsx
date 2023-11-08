import React from "react";
import NextLink from "next/link";
import "bootstrap/dist/css/bootstrap.min.css";
import { ToastContainer, toast } from "react-toastify";
import { ToastContext } from "lib/useMessage";
import "react-toastify/dist/ReactToastify.css";

const toastProvider = {
  info: (msg: string) => {
    toast.info(msg, {
      autoClose: 2000,
      closeOnClick: true,
    });
  },
  error: (msg: string) => {
    toast.error(msg, {
      autoClose: 2000,
      closeOnClick: true,
    });
  },
};

const AdminLayout = (page: React.ReactNode) => {
  return (
    <>
      <ToastContext.Provider value={toastProvider}>
        <div className="container">
          <nav className="navbar navbar-expand-lg bg-light">
            <ul className="navbar-nav">
              <li className="nav-item">
                <NextLink
                  href="/admin"
                  className="nav-link active"
                  aria-current="page"
                >
                  Config
                </NextLink>
              </li>

              <li className="nav-item">
                <NextLink
                  href="/admin/operations"
                  className="nav-link active"
                  aria-current="page"
                >
                  operations
                </NextLink>
              </li>

              <li className="nav-item">
                <NextLink
                  href="/admin/teams"
                  className="nav-link active"
                  aria-current="page"
                >
                  teams
                </NextLink>
              </li>

              <li className="nav-item">
                <NextLink
                  href="/admin/tasks"
                  className="nav-link active"
                  aria-current="page"
                >
                  tasks
                </NextLink>
              </li>
            </ul>
          </nav>
          {page}
        </div>
        <ToastContainer />
      </ToastContext.Provider>
    </>
  );
};

export default AdminLayout;
