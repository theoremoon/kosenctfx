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
                <NextLink href="/admin">
                  <a className="nav-link active" aria-current="page">
                    Config
                  </a>
                </NextLink>
              </li>

              <li className="nav-item">
                <NextLink href="/admin/operations">
                  <a className="nav-link active" aria-current="page">
                    operations
                  </a>
                </NextLink>
              </li>

              <li className="nav-item">
                <NextLink href="/admin/teams">
                  <a className="nav-link active" aria-current="page">
                    teams
                  </a>
                </NextLink>
              </li>

              <li className="nav-item">
                <NextLink href="/admin/tasks">
                  <a className="nav-link active" aria-current="page">
                    tasks
                  </a>
                </NextLink>
              </li>

              <li className="nav-item">
                <NextLink href="/admin/agents">
                  <a className="nav-link active" aria-current="page">
                    agents
                  </a>
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
