<template>
    <div class="container">
        <div>
            <Card dis-hover>
                <p slot="title">
                    <Icon type="ios-build-outline"></Icon>
                    ES日志查询
                    <span style="float: right; margin-top: -2px; margin-right: 20px">
                        <Badge :count="server" type="primary"></Badge>
                        <span> </span>
                        <Badge :count="hit" type="success"></Badge>
                    </span>
                </p>
                <ul>
                    <Form ref="formField" :model="queryForm" :rules="ruleValidate" :label-width="80">
                        <FormItem label="ES地址" prop="url">
                            <Input v-model="queryForm.url" placeholder="ES地址" type="text" :rows="4"/>
                        </FormItem>
                        <FormItem label="账户">
                            <Input v-model="queryForm.auth_user" placeholder="ES认证账户" type="text" :rows="4"/>
                        </FormItem>
                        <FormItem label="密码">
                            <Input v-model="queryForm.auth_password" placeholder="ES认证密码" type="password" :rows="4"/>
                        </FormItem>

                        <Divider>查询条件</Divider>

                        <FormItem label="索引" prop="index">
                            <Input v-model="queryForm.index" placeholder="ES索引" type="text" :rows="4"/>
                        </FormItem>
                        <FormItem label="条件" prop="json">
                            <Input v-model="queryForm.json" placeholder="匹配条件,json格式" type="textarea" :rows="4"/>
                        </FormItem>
                        <FormItem label="时间范围" prop="date_range">
                            <Row>
                                <FormItem prop="date_range">
                                    <DatePicker type="daterange" placeholder="日期" v-model="queryForm.date_range" style="width: 206px"/>
                                </FormItem>
                            </Row>
                        </FormItem>
                    </Form>

                    <Form :label-width="80">
                        <FormItem>
                            <div style="float: right">
                                <Button type="primary" @click="submit">查询</Button>
                                <Button type="error" style="margin-left: 8px" @click="download">导出结果</Button>
                            </div>
                        </FormItem>
                    </Form>
                </ul>
            </Card>
        </div>
    </div>
</template>

<script>
    export default {
        name: "Home",
        data() {
            return {
                hit: "HIT数量: 0",
                server: "ES服务状态: 未连接",
                queryForm: {
                    url: '',
                    auth_user: '',
                    auth_password: '',
                    index: '',
                    date_range: [],
                    json: '',
                    query: {},
                },
                formField: {
                    url: '',
                    date_range: [],
                    index: '',
                    json: '',
                },
                downloadObj: {},
                ruleValidate: {
                    url: [
                        { required: true, message: '请设置ES查询地址', trigger: 'blur' }
                    ],
                    date_range: [
                        { required: true, type: 'array', message: '请选择查询时间', trigger: 'blur'}
                    ],
                    index: [
                        { required: true, type: 'string', message: '请设置索引', trigger: 'blur' }
                    ],
                    json: [
                        { required: true, message: '请设置查询条件', trigger: 'blur' }
                    ]
                }
            }
        },
        methods: {
            submit() {
                if (this.queryForm.url == '') {
                    this.$Message.warning('请设置ES查询地址')
                    return
                }
                if (this.queryForm.date == '') {
                    this.$Message.warning('请设置正确查询时间')
                    return
                }

                this.queryForm.query.index = this.queryForm.index ;
                this.queryForm.query.match = JSON.parse(this.queryForm.json);


                let sDate = this.queryForm.date_range[0];
                let eDate = this.queryForm.date_range[1];

                // fix timezone
                const offset = sDate.getTimezoneOffset();
                sDate = new Date(sDate.getTime() - (offset*60*1000));
                eDate = new Date(eDate.getTime() - (offset*60*1000));

                this.queryForm.query.start_date = sDate.toISOString().split('T')[0];
                this.queryForm.query.end_date = eDate.toISOString().split('T')[0];

                console.log("form: ", this.queryForm.query);

                let self = this;
                window.backend.App.Search(this.queryForm).then(function(r) {
                    console.log(r)

                    if (r.code === 0) {
                        self.downloadObj = self.queryForm;
                        self.server = "ES服务状态: ok";
                        self.hit = "HIT数量: " + r.data.Hit +",耗时:" + r.data.CostTime;

                    } else {
                        self.server = "ES服务状态: 异常";
                        self.$Message.warning({ content: r.message, duration: 5, closable: true});
                    }

                    // this.$Loading.finish();
                })
            },
            download() {
                if (this.downloadObj === "") {
                    this.$Message.warning("请先执行查询操作"); return;
                }

                this.$Loading.start();

                let self = this;
                window.backend.App.Download(this.downloadObj).then(function(r) {
                    console.log(r)
                    if (r.code === 0) {
                        self.$Message.success("导出数据成功");
                        self.$Message.success({content: "文件保存路径: " + r.data.Data.filename, duration: 10, closable: true});
                    } else {
                        self.$Message.warning({ content: r.message, duration: 5, closable: true});
                    }

                    self.$Loading.finish();
                })
            }
        }
    }
</script>

<style scoped>
</style>
