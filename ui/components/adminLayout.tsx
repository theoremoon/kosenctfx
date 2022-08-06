import React from "react";
import NextLink from "next/link";
import "bootstrap/dist/css/bootstrap.min.css";

const AdminLayout = (page: React.ReactNode) => {
  return (
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
        </ul>
      </nav>
      {page}
    </div>
  );
};

export default AdminLayout;
