<template>
    <table>
        <thead>
            <tr>Rank</tr>
            <tr>Team</tr>
        </thead>
        <tbody></tbody>
    </table>
</template>

<script>
export default {
    data() {
        return {
            standings: [
                {
                    pos: 1,
                    team: "zer0pts",
                    id: 3,
                    score: 1900,
                    solved_tasks: [
                        {
                            name: "padrsa",
                            time: 1589091200,
                            score: 1000,
                            solved_by: "ptr-yudai",
                        },
                        {
                            name: "padrsa2",
                            time: 1589091100,
                            score: 900,
                            solved_by: "theoremoon",
                        },
                    ],
                },
                {
                    pos: 2,
                    team: "TokyoWesterns",
                    id: 5,
                    score: 1900,
                    solved_tasks: [
                        {
                            name: "padrsa",
                            time: 1589091230,
                            score: 1000,
                            solved_by: "t0nk42",
                        },
                        {
                            name: "padrsa2",
                            time: 1589091110,
                            score: 900,
                            solved_by: "N4NU",
                        },
                    ],
                },
            ]
        }
    },
}
</script>