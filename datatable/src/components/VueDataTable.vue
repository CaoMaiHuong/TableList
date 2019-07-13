<template>
  <div class="animated fadeIn">
      <b-row>
          <b-col sm="12" lg="12">
              <b-card title header-tag="header" footer-tag="footer">
                    <b-col md="6" class="my-1">
                        <b-form-group label-cols-sm="3" label="Filter" class="mb-0">
                        <b-input-group>
                            <b-form-input v-model="filter" placeholder="Type to Search"></b-form-input>
                            <b-input-group-append>
                            <b-button :disabled="!filter" @click="filter = ''">Clear</b-button>
                            </b-input-group-append>
                        </b-input-group>
                        </b-form-group>
                    </b-col>
                    <b-table 
                        ref="table"
                        :busy.sync="isBusy"
                        :items="myProvider"
                        :fields="fields"
                        :current-page="currentPage"
                        :per-page="perPage"
                        striped
                        bordered
                        small
                        primary-key="id"
                        :filter="filter"
                    >
                    </b-table>
              </b-card>
          </b-col>
      </b-row>
      <b-row style="float: right">
        <b-col md="6">
            <b-pagination
            v-model="currentPage"
            :total-rows="totalRows"
            :per-page="perPage"
            first-text="<<"
            prev-text="<"
            next-text=">"
            last-text=">>"
            >   
            </b-pagination>
        </b-col>
      </b-row>
  </div>
</template>

<script>
    export default{
        name: "tablelist",
        data(){
            return {
                items: [],
                fields: [
                {
                    key: 'name',
                    label: 'Name',
                    sortable: true
                },
                {
                    key: 'family',
                    label: 'Family',
                    sortable: true
                },
                {
                    key: 'cvss_base',
                    label: 'Severity',
                    sortable: true
                },
                {
                    key: 'creation_time',
                    label: 'Created',
                    sortable: true
                },
                {
                    key: 'qod',
                    label: 'Qod',
                    sortable: true
                }
                ],
                isBusy: false,
                totalRows:1,
                currentPage:1,
                perPage:15,
                filter:'',
                sort:'asc',
                order:'id'
            }
        },
        methods: {
            myProvider(ctx) {
                if (ctx.filter != "") {
                    var items = []
                    this.items.filter((value,index) => {
                        if(value.name.toLowerCase().indexOf(ctx.filter.toLowerCase()) > -1 || value.family.toLowerCase().indexOf(ctx.filter.toLowerCase()) > -1 || value.cvss_base.toString().toLowerCase().indexOf(ctx.filter.toLowerCase()) > -1 || value.qod.toString().toLowerCase().indexOf(ctx.filter.toLowerCase()) > -1) {
                            items.push(value)
                        }
                    })
                    if(items.length > 0) {
                        return items
                    } else {
                        return []
                    }
                } else {
                    this.isBusy = true
                    if (ctx.sortBy) {
                        let url1 = "http://localhost:8081/nvts/page=" + ctx.currentPage
                        let url2 = `${url1}&_order=${ctx.sortBy } ${ctx.sortDesc ? 'desc' : 'asc'}`
                        let promise = axios.get(url2)
                        return promise.then(res => {
                            var items = res.data.datas
                            this.items = res.data.datas
                            this.currentPage = res.data.current_page
                            this.totalRows = res.data.total
                            this.isBusy = false
                            return items
                        })
                    } else {
                        let url = "http://localhost:8081/nvts/page=" + ctx.currentPage + "&_order=" + this.order +" "+this.sort
                        let promise = axios.get(url)
                    
                        return promise.then(res => {
                            var items = res.data.datas
                            this.items = res.data.datas
                            this.currentPage = res.data.current_page
                            this.totalRows = res.data.total
                            this.isBusy = false
                            return items
                        }).catch(error => {
                            this.isBusy = false
                            return []
                        })
                    }
                }
            }
        }
    }
</script>

<style>
#app {
    width: 92%;
    margin: 55px;
}
.input-group > .input-group-append > .btn{
    font-size: 14px;
}

body{
    font-size: 14px;
}
.VuePagination {
  text-align: center;
}

.vue-title {
  text-align: center;
  margin-bottom: 10px;
}

.vue-pagination-ad {
  text-align: center;
}
.form-control{
    font-size: 14px;
}
.table>tbody>tr>td{
    text-align: left;
}

.VueTables__search-field{
    display: flex;
    margin-bottom: 20px;
}
.table>thead:first-child>tr:first-child>th{
    border-top: 1px solid #ddd;
    vertical-align: middle;
}
.glyphicon.glyphicon-eye-open {
  width: 16px;
  display: block;
  margin: 0 auto;
}
.form-inline label{
    justify-content: left;
    margin-right: 12px;
}

th:nth-child(3) {
  text-align: center;
}

.VueTables__child-row-toggler {
  width: 16px;
  height: 16px;
  line-height: 16px;
  display: block;
  margin: auto;
  text-align: center;
}

.VueTables__child-row-toggler--closed::before {
  content: "+";
}

.VueTables__child-row-toggler--open::before {
  content: "-";
}
.glyphicon {
    font-size: 13px;
    top: 4px;
}
.form-row{
    margin-bottom: 10px;
}
legend{
    border: none;
    font-weight: bold;
}
</style>