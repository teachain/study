​         <FormItem label="部署环境:" :label-width="100" prop="environment_id">

​          <Select v-model="formItem.environment_id">

​            <Option v-for="(item,index) in project_environments" :key="'environment'+index" :value="item.id" :label="item.name">

​            </Option>

​          </Select>

​        </FormItem>
  在代码里操作project_environments，不要随便操作"formItem.environment_id，比如说你可以清空project_environments这个数组，但是不要去设置"formItem.environment_id=0