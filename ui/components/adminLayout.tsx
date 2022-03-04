import {
  Box,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Stack,
} from "@chakra-ui/react";
import React from "react";
import NextLink from "next/link";

interface AdminLayoutProps {
  children: React.ReactNode;
}

const AdminLayout = ({ children, ...props }: AdminLayoutProps) => {
  return (
    <Box mt="5">
      <Breadcrumb>
        <BreadcrumbItem>
          <NextLink href="/admin" passHref>
            <BreadcrumbLink>Config</BreadcrumbLink>
          </NextLink>
        </BreadcrumbItem>

        <BreadcrumbItem>
          <NextLink href="/admin/tasks" passHref>
            <BreadcrumbLink>Tasks</BreadcrumbLink>
          </NextLink>
        </BreadcrumbItem>

        <BreadcrumbItem>
          <NextLink href="/admin/teams" passHref>
            <BreadcrumbLink>Teams</BreadcrumbLink>
          </NextLink>
        </BreadcrumbItem>
      </Breadcrumb>

      <Stack mt={5}>{children}</Stack>
    </Box>
  );
};

export default AdminLayout;
