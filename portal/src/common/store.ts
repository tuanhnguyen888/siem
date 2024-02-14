import { createStore } from "vuex";
import router from "./router";
import { ICurrentPage } from "./model";
import { demoEvent } from "../components/demo-data";
interface IState {
  isLogin: boolean;
  currentPage: string;
  tableDatas: any[];
}

const store = createStore<IState>({
  state: {
    isLogin: false,
    currentPage: ICurrentPage.Event,
    tableDatas: demoEvent.slice(0, 7),
  },
  mutations: {
    setCurrentPage(state, payload) {
      state.currentPage = payload;
    },
    setLogin(state, payload) {
      state.isLogin = payload;
    },
    setTableDatas(state, payload) {
      state.tableDatas = state.tableDatas.concat(payload);
    },
  },
  actions: {
    getLogin(context) {
      context.commit("setLogin", true);
      router.push({ name: "event" });
    },
    getMoreTableDatas(context) {
      if (context.state.tableDatas > demoEvent.length) {
        return;
      }
      const num = context.state.tableDatas;
      const arr = demoEvent.slice(num, num + 5);

      context.commit("setTableDatas", arr);
    },
    getCurrentPage(context, payload) {
      context.commit("setCurrentPage", payload);
    },
  },
});

export default store;
