import type {Issue} from "@/api/issues/types.ts";


export const useIssueDialogStore = defineStore('issue-dialog', {
  state: () => ({
    isOpen: false,
    issue: null as Issue | null,
  }),
  actions: {
    open(issue: Issue) {
      this.isOpen = true;
      this.issue = issue;
      console.log("Issue dialog opened for issue:", issue)
    },
    close() {
      this.isOpen = false;
      this.issue = null;
    },
  },
});
