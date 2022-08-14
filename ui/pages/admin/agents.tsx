import AdminLayout from "components/adminLayout";
import useAgents, { Agent, fetchAgents } from "lib/api/admin/agent";
import { dateFormat } from "lib/date";
import { GetStaticProps } from "next";

const Agents = () => {
  const { data: agents } = useAgents();

  return (
    <>
      <h5 className="mt-4">challenge instances</h5>
      <table className="table">
        <thead>
          <tr>
            <th>ID</th>
            <th>IP Address</th>
            <th>Last Activity</th>
          </tr>
        </thead>
        <tbody>
          {agents &&
            agents.map((a) => (
              <tr key={a.agent_id}>
                <td>{a.agent_id}</td>
                <td>{a.public_ip}</td>
                <td>{dateFormat(a.last_activity_at)}</td>
              </tr>
            ))}
        </tbody>
      </table>
    </>
  );
};

export const getStaticProps: GetStaticProps = async () => {
  return {
    props: {},
  };
};

Agents.getLayout = AdminLayout;

export default Agents;
