import { Task } from "lib/api/tasks";
import { useRouter } from "next/router";
import React from "react";
import styles from "./taskCard.module.scss";
import cx from "classnames";

type TaskCardProps = {
  task: Task;
  isSolved: boolean;
} & React.ComponentPropsWithoutRef<"div">;

const TaskCard = ({ task, isSolved }: TaskCardProps) => {
  const router = useRouter();
  const cls = isSolved
    ? cx(styles["task"], styles["task-solved"])
    : styles["task"];
  return (
    <a
      className={cls}
      onClick={() => {
        router.push("/tasks/" + task.id, undefined, {
          scroll: false,
          shallow: true,
        });
      }}
    >
      <div className={styles["task-cover"]}>
        <div className={styles["task-upper"]}>
          <div className={styles["task-name"]}>{task.name}</div>
        </div>
        <div className={styles["task-middle"]}>
          <div className={styles["task-score"]}>{task.score}</div>
        </div>
        <div className={styles["task-lower"]}>
          <div className={styles["task-tags"]}>
            {task.tags.map((tag) => (
              <div className={styles["task-tag"]} key={tag}>
                {tag}
              </div>
            ))}
          </div>

          <div className={styles["task-solves"]}>
            {task.solved_by.length} solves
          </div>
        </div>
      </div>
    </a>
  );
};

export default TaskCard;
