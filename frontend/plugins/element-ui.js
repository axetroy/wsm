import Vue from 'vue'

import {
  Row,
  Col,
  Dropdown,
  DropdownMenu,
  DropdownItem,
  Popover,
  Form,
  FormItem,
  Input,
  InputNumber,
  Button,
  Container,
  Menu,
  MenuItem,
  MenuItemGroup,
  Submenu,
  Aside,
  Header,
  Main,
  Table,
  TableColumn,
  Breadcrumb,
  BreadcrumbItem,
  Upload,
  DatePicker,
  TimePicker,
  Select,
  Option,
  Pagination,
  Dialog,
  Radio,
  RadioGroup,
  Alert,
  Switch,
  Tag,
  CheckboxGroup,
  Checkbox,
  Message,
  MessageBox,
  Tooltip,
  Card,
  Tabs,
  TabPane,
  Loading,
  Popconfirm,
  Autocomplete
} from 'element-ui'

// message
Vue.prototype.$message = Vue.prototype.$message || Message
Vue.prototype.$success = Vue.prototype.$success || Message.success.bind(Message)
Vue.prototype.$info = Vue.prototype.$info || Message.info.bind(Message)
Vue.prototype.$warning = Vue.prototype.$warning || Message.warning.bind(Message)
Vue.prototype.$error = Vue.prototype.$error || Message.error.bind(Message)

// messageBox
Vue.prototype.$alert = Vue.prototype.$alert || MessageBox.alert.bind(MessageBox)
Vue.prototype.$confirm =
  Vue.prototype.$confirm || MessageBox.confirm.bind(MessageBox)
Vue.prototype.$prompt =
  Vue.prototype.$prompt || MessageBox.prompt.bind(MessageBox)

Vue.use(Row)
Vue.use(Col)
Vue.use(Dropdown)
Vue.use(DropdownMenu)
Vue.use(DropdownItem)
Vue.use(Popover)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Input)
Vue.use(InputNumber)
Vue.use(Button)
Vue.use(Container)
Vue.use(Menu)
Vue.use(MenuItem)
Vue.use(MenuItemGroup)
Vue.use(Submenu)
Vue.use(Aside)
Vue.use(Header)
Vue.use(Main)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Breadcrumb)
Vue.use(BreadcrumbItem)
Vue.use(Upload)
Vue.use(DatePicker)
Vue.use(TimePicker)
Vue.use(Select)
Vue.use(Option)
Vue.use(Pagination)
Vue.use(Dialog)
Vue.use(Radio)
Vue.use(RadioGroup)
Vue.use(Alert)
Vue.use(Switch)
Vue.use(Tag)
Vue.use(CheckboxGroup)
Vue.use(Checkbox)
Vue.use(Tooltip)
Vue.use(Card)
Vue.use(Tabs)
Vue.use(TabPane)
Vue.use(Popconfirm)
Vue.use(Autocomplete)

Vue.use(Loading.directive)
Vue.prototype.$loading = Loading.service
