window.BENCHMARK_DATA = {
  "lastUpdate": 1733731421259,
  "repoUrl": "https://github.com/sunerpy/requests",
  "entries": {
    "Requests-Benchmark": [
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "11db1ffcf1d4648b3c64e8a4b7f9e74e23cf9c62",
          "message": "ci(benchmark): 更新性能基准测试工作流\n\n- 使用 benchmark-action 生成和存储基准测试结果\n- 自动推送基准测试报告到 GitHub Pages",
          "timestamp": "2024-12-08T12:35:12+08:00",
          "tree_id": "58c433b8892c0c60893cd67a0b3fda08110a6947",
          "url": "https://github.com/sunerpy/requests/commit/11db1ffcf1d4648b3c64e8a4b7f9e74e23cf9c62"
        },
        "date": 1733632581001,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 97402,
            "unit": "ns/op",
            "extra": "12364 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100970,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100899,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100560,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25397,
            "unit": "ns/op",
            "extra": "47257 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20508,
            "unit": "ns/op",
            "extra": "58208 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 24693,
            "unit": "ns/op",
            "extra": "58712 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20981,
            "unit": "ns/op",
            "extra": "57578 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "1f5a5f35c12639fc753627e55a5ed1e395ae395e",
          "message": "docs: 更新 README 文件并添加中文版本\n\n- 更新英文 README.md 文件内容\n- 新增中文 readme-cn.md 文件\n- 优化文档结构，增加目录和徽章\n- 更新示例代码和 API 文档",
          "timestamp": "2024-12-08T13:24:58+08:00",
          "tree_id": "6ca69b4379cc17bcd5eca42c26ece664d5426412",
          "url": "https://github.com/sunerpy/requests/commit/1f5a5f35c12639fc753627e55a5ed1e395ae395e"
        },
        "date": 1733635555406,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94744,
            "unit": "ns/op",
            "extra": "12327 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98135,
            "unit": "ns/op",
            "extra": "12172 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98580,
            "unit": "ns/op",
            "extra": "12192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 97966,
            "unit": "ns/op",
            "extra": "12200 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 27505,
            "unit": "ns/op",
            "extra": "48156 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20632,
            "unit": "ns/op",
            "extra": "55200 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20624,
            "unit": "ns/op",
            "extra": "58003 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20922,
            "unit": "ns/op",
            "extra": "58866 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "7ea5ab731a6c5432699c5f151dbcd6cb70f0c1ea",
          "message": "docs: 更新 README 文件并添加中文版本\n\n- 更新英文 README.md 文件内容\n- 新增中文 readme-cn.md 文件\n- 优化文档结构，增加目录和徽章\n- 更新示例代码和 API 文档\n- 新增 Benchmark 工作流，运行性能测试生成 HTML 报告并上传到 GitHub Pages",
          "timestamp": "2024-12-08T13:28:58+08:00",
          "tree_id": "6ca69b4379cc17bcd5eca42c26ece664d5426412",
          "url": "https://github.com/sunerpy/requests/commit/7ea5ab731a6c5432699c5f151dbcd6cb70f0c1ea"
        },
        "date": 1733635869599,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95525,
            "unit": "ns/op",
            "extra": "12519 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 99849,
            "unit": "ns/op",
            "extra": "12178 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98913,
            "unit": "ns/op",
            "extra": "12192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98642,
            "unit": "ns/op",
            "extra": "12213 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24988,
            "unit": "ns/op",
            "extra": "48128 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20210,
            "unit": "ns/op",
            "extra": "59992 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20913,
            "unit": "ns/op",
            "extra": "60295 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20705,
            "unit": "ns/op",
            "extra": "58940 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "c8e51a9caf2a6eff1cb52238a31a7a5f31b1b9c9",
          "message": "feat(requests): 添加自定义 DNS 服务器支持\n\n- 在 Session 类型中增加 WithDNS 方法，支持设置自定义 DNS 服务器\n- 实现自定义 DNS 解析逻辑\n\n- 增加相关单元测试，验证自定义 DNS 功能的正确性",
          "timestamp": "2024-12-08T20:52:17+08:00",
          "tree_id": "55f8c734827a3da5f7ef356e94402bbbb2a1c2d6",
          "url": "https://github.com/sunerpy/requests/commit/c8e51a9caf2a6eff1cb52238a31a7a5f31b1b9c9"
        },
        "date": 1733662384100,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95805,
            "unit": "ns/op",
            "extra": "12405 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 99465,
            "unit": "ns/op",
            "extra": "12058 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100949,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100439,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25286,
            "unit": "ns/op",
            "extra": "46009 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20489,
            "unit": "ns/op",
            "extra": "58611 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20552,
            "unit": "ns/op",
            "extra": "58785 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20754,
            "unit": "ns/op",
            "extra": "59725 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "1389d0675c25752a642efe3278fbbdc0ea0b9fc8",
          "message": "feat(requests): 添加自定义 DNS 服务器支持和session重置\n\n- 在 Session 类型中增加 WithDNS 方法，支持设置自定义 DNS 服务器\n- 实现自定义 DNS 解析逻辑\n\n- 增加相关单元测试，验证自定义 DNS 功能的正确性\n\n- 增加了 Clear() 方法，确保会话状态重置",
          "timestamp": "2024-12-08T21:02:49+08:00",
          "tree_id": "848bb5b4250328645ac166e78ef396050097e0bc",
          "url": "https://github.com/sunerpy/requests/commit/1389d0675c25752a642efe3278fbbdc0ea0b9fc8"
        },
        "date": 1733663076655,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94821,
            "unit": "ns/op",
            "extra": "12531 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 97851,
            "unit": "ns/op",
            "extra": "12192 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98365,
            "unit": "ns/op",
            "extra": "12158 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98077,
            "unit": "ns/op",
            "extra": "12205 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25751,
            "unit": "ns/op",
            "extra": "47043 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 23659,
            "unit": "ns/op",
            "extra": "45820 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20343,
            "unit": "ns/op",
            "extra": "57511 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20751,
            "unit": "ns/op",
            "extra": "59461 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "b4f9c20a400738d5719449b6bf8979608c50c3df",
          "message": "ci: 更新 GitHub Actions 工作流触发条件\n\n- 在 benchmark.yml 和 release.yml 中添加 dev 分支的推送触发\n- 在 benchmark.yml 中添加 main 分支的 PR 触发",
          "timestamp": "2024-12-09T11:45:54+08:00",
          "tree_id": "b1fd00253b2a1bcba1edde328c7e71cc62f2ab82",
          "url": "https://github.com/sunerpy/requests/commit/b4f9c20a400738d5719449b6bf8979608c50c3df"
        },
        "date": 1733716023147,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94034,
            "unit": "ns/op",
            "extra": "12760 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98214,
            "unit": "ns/op",
            "extra": "12235 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 99346,
            "unit": "ns/op",
            "extra": "12130 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98590,
            "unit": "ns/op",
            "extra": "12201 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 30015,
            "unit": "ns/op",
            "extra": "48313 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20370,
            "unit": "ns/op",
            "extra": "54230 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20655,
            "unit": "ns/op",
            "extra": "58539 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20634,
            "unit": "ns/op",
            "extra": "59032 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "b4f9c20a400738d5719449b6bf8979608c50c3df",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/b4f9c20a400738d5719449b6bf8979608c50c3df"
        },
        "date": 1733716157995,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95758,
            "unit": "ns/op",
            "extra": "12364 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100070,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100494,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 99559,
            "unit": "ns/op",
            "extra": "12010 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25265,
            "unit": "ns/op",
            "extra": "46749 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 23320,
            "unit": "ns/op",
            "extra": "58154 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20375,
            "unit": "ns/op",
            "extra": "57363 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20264,
            "unit": "ns/op",
            "extra": "59478 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "1a67006cc63ded39651f252807934bb674e47885",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/1a67006cc63ded39651f252807934bb674e47885"
        },
        "date": 1733721919738,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94590,
            "unit": "ns/op",
            "extra": "12471 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98893,
            "unit": "ns/op",
            "extra": "12566 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 99422,
            "unit": "ns/op",
            "extra": "12141 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100184,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24943,
            "unit": "ns/op",
            "extra": "48746 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 21030,
            "unit": "ns/op",
            "extra": "59098 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20916,
            "unit": "ns/op",
            "extra": "58503 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20651,
            "unit": "ns/op",
            "extra": "56932 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "d5660fcabb1a70b2adf3eb6e540a34d18925d849",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/d5660fcabb1a70b2adf3eb6e540a34d18925d849"
        },
        "date": 1733722063494,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 97384,
            "unit": "ns/op",
            "extra": "12174 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100355,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 101188,
            "unit": "ns/op",
            "extra": "12022 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 101007,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24926,
            "unit": "ns/op",
            "extra": "46803 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20095,
            "unit": "ns/op",
            "extra": "58850 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20420,
            "unit": "ns/op",
            "extra": "58843 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 21216,
            "unit": "ns/op",
            "extra": "59574 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "1bbd1d863a1a3ca50db42746b55919829b494c99",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/1bbd1d863a1a3ca50db42746b55919829b494c99"
        },
        "date": 1733724731906,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95481,
            "unit": "ns/op",
            "extra": "12428 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98892,
            "unit": "ns/op",
            "extra": "12115 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98527,
            "unit": "ns/op",
            "extra": "12123 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98979,
            "unit": "ns/op",
            "extra": "12097 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 27281,
            "unit": "ns/op",
            "extra": "47341 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 22143,
            "unit": "ns/op",
            "extra": "45915 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20683,
            "unit": "ns/op",
            "extra": "58839 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20113,
            "unit": "ns/op",
            "extra": "59924 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "ac796206a5b9ed0c2cf6ca59cb8fce5556b8b8ef",
          "message": "添加 WithBasicAuth",
          "timestamp": "2024-12-08T13:09:20Z",
          "url": "https://github.com/sunerpy/requests/pull/1/commits/ac796206a5b9ed0c2cf6ca59cb8fce5556b8b8ef"
        },
        "date": 1733725500162,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95797,
            "unit": "ns/op",
            "extra": "12278 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 99736,
            "unit": "ns/op",
            "extra": "12024 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100018,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 99583,
            "unit": "ns/op",
            "extra": "12014 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24835,
            "unit": "ns/op",
            "extra": "47808 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 24998,
            "unit": "ns/op",
            "extra": "57620 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20492,
            "unit": "ns/op",
            "extra": "58162 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20636,
            "unit": "ns/op",
            "extra": "57961 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "695a40fffbf6e22b862e8cec5ec6ba297d648b19",
          "message": "添加 WithBasicAuth (#1)\n\n* feat(requests): 添加基本认证功能\r\n\r\n- 在 Session 接口中添加 WithBasicAuth 方法\r\n- 实现基本认证逻辑，生成 Base64 编码的认证头\r\n- 更新单元测试，增加对基本认证功能的测试\r\n\r\n\r\n* ci: 更新工作流以包含代码覆盖率检查和自动审批\r\n\r\n- 在 go-test.yml 中添加代码覆盖率检查步骤\r\n- 如果覆盖率低于 90%，则失败\r\n- 新增 Auto Approve PR 工作流，自动审批通过测试的 PR",
          "timestamp": "2024-12-09T14:27:47+08:00",
          "tree_id": "b6e2b8b3a7e25c0853bc7285a42a844988b64627",
          "url": "https://github.com/sunerpy/requests/commit/695a40fffbf6e22b862e8cec5ec6ba297d648b19"
        },
        "date": 1733725707043,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 95237,
            "unit": "ns/op",
            "extra": "12730 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 92775,
            "unit": "ns/op",
            "extra": "12168 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98123,
            "unit": "ns/op",
            "extra": "12142 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98887,
            "unit": "ns/op",
            "extra": "12196 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 24972,
            "unit": "ns/op",
            "extra": "47208 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 21267,
            "unit": "ns/op",
            "extra": "54360 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20216,
            "unit": "ns/op",
            "extra": "57548 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20762,
            "unit": "ns/op",
            "extra": "57081 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "9f1ab2d80110a2e880758d47529f1c4f19be1fb8",
          "message": "ci: 移除 auto-approve 工作流中的依赖任务",
          "timestamp": "2024-12-09T06:27:51Z",
          "url": "https://github.com/sunerpy/requests/pull/2/commits/9f1ab2d80110a2e880758d47529f1c4f19be1fb8"
        },
        "date": 1733726387397,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 96366,
            "unit": "ns/op",
            "extra": "12283 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 99424,
            "unit": "ns/op",
            "extra": "12075 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100315,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 101207,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25597,
            "unit": "ns/op",
            "extra": "44066 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20337,
            "unit": "ns/op",
            "extra": "58186 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20700,
            "unit": "ns/op",
            "extra": "57674 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20725,
            "unit": "ns/op",
            "extra": "59478 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "distinct": true,
          "id": "c3959180888e0529a6f5e766811f631a40985342",
          "message": "ci: 移除 auto-approve 工作流中的依赖任务\n\n- 删除了 approval job 中的 needs 字段",
          "timestamp": "2024-12-09T14:47:05+08:00",
          "tree_id": "fce07c2cbe25ac127059ea042e898835fc003e19",
          "url": "https://github.com/sunerpy/requests/commit/c3959180888e0529a6f5e766811f631a40985342"
        },
        "date": 1733726857179,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 97099,
            "unit": "ns/op",
            "extra": "12229 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100758,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 101065,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100907,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 29449,
            "unit": "ns/op",
            "extra": "41865 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20392,
            "unit": "ns/op",
            "extra": "58692 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20559,
            "unit": "ns/op",
            "extra": "57391 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20611,
            "unit": "ns/op",
            "extra": "58622 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "0b5f6d9af4acc6500c33c63b930f1384bbc5292e",
          "message": "Dev",
          "timestamp": "2024-12-09T06:47:10Z",
          "url": "https://github.com/sunerpy/requests/pull/3/commits/0b5f6d9af4acc6500c33c63b930f1384bbc5292e"
        },
        "date": 1733727508128,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 97557,
            "unit": "ns/op",
            "extra": "12010 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100179,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 99869,
            "unit": "ns/op",
            "extra": "12092 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 100257,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 26220,
            "unit": "ns/op",
            "extra": "45938 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20678,
            "unit": "ns/op",
            "extra": "58200 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 24642,
            "unit": "ns/op",
            "extra": "58191 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20525,
            "unit": "ns/op",
            "extra": "59306 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "802733cd4227431f70dc4e7511a3619035ac27c8",
          "message": "Dev",
          "timestamp": "2024-12-09T06:47:10Z",
          "url": "https://github.com/sunerpy/requests/pull/3/commits/802733cd4227431f70dc4e7511a3619035ac27c8"
        },
        "date": 1733727543542,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 96641,
            "unit": "ns/op",
            "extra": "12244 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 99614,
            "unit": "ns/op",
            "extra": "12009 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98574,
            "unit": "ns/op",
            "extra": "12135 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 99999,
            "unit": "ns/op",
            "extra": "12086 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25024,
            "unit": "ns/op",
            "extra": "45811 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20206,
            "unit": "ns/op",
            "extra": "58282 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20664,
            "unit": "ns/op",
            "extra": "58321 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20769,
            "unit": "ns/op",
            "extra": "55996 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ce57d46a173737d1b425a3bbf7b3efd9cba75a66",
          "message": "Dev (#3)\n\n* ci: 移除 auto-approve 工作流中的依赖任务\r\n\r\n- 删除了 approval job 中的 needs 字段\r\n- 添加检查，跳过非 pull request 触发的工作流程\r\n- 为检查 PR 事件步骤添加 id，以便于后续引用",
          "timestamp": "2024-12-09T15:01:15+08:00",
          "tree_id": "477172a0c7d6532957267a0c9120b75bac0ca4ce",
          "url": "https://github.com/sunerpy/requests/commit/ce57d46a173737d1b425a3bbf7b3efd9cba75a66"
        },
        "date": 1733727710312,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94199,
            "unit": "ns/op",
            "extra": "12525 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 97680,
            "unit": "ns/op",
            "extra": "12338 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 97754,
            "unit": "ns/op",
            "extra": "12266 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98669,
            "unit": "ns/op",
            "extra": "12243 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 26265,
            "unit": "ns/op",
            "extra": "46410 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 23139,
            "unit": "ns/op",
            "extra": "45184 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20546,
            "unit": "ns/op",
            "extra": "57506 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20514,
            "unit": "ns/op",
            "extra": "58026 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "417398a46768a592e87fe9dd71655bb1de94b2bc",
          "message": "Dev",
          "timestamp": "2024-12-09T07:01:19Z",
          "url": "https://github.com/sunerpy/requests/pull/4/commits/417398a46768a592e87fe9dd71655bb1de94b2bc"
        },
        "date": 1733730541583,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 96872,
            "unit": "ns/op",
            "extra": "12495 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100374,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100301,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 85687,
            "unit": "ns/op",
            "extra": "12074 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 26060,
            "unit": "ns/op",
            "extra": "41606 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20474,
            "unit": "ns/op",
            "extra": "57394 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20371,
            "unit": "ns/op",
            "extra": "57753 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20698,
            "unit": "ns/op",
            "extra": "58501 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "49fb4819f848db06ce8e37e55fdef6f4a95a705b",
          "message": "Dev (#4)\n\n* ci: 移除 auto-approve 工作流中的依赖任务\r\n\r\n- 删除了 approval job 中的 needs 字段",
          "timestamp": "2024-12-09T15:50:13+08:00",
          "tree_id": "24722565f2b86e8f41ed93019cf5f6f3a13addac",
          "url": "https://github.com/sunerpy/requests/commit/49fb4819f848db06ce8e37e55fdef6f4a95a705b"
        },
        "date": 1733730651730,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 96092,
            "unit": "ns/op",
            "extra": "12404 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100647,
            "unit": "ns/op",
            "extra": "12123 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 100646,
            "unit": "ns/op",
            "extra": "12100 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 99908,
            "unit": "ns/op",
            "extra": "12010 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25067,
            "unit": "ns/op",
            "extra": "48510 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20484,
            "unit": "ns/op",
            "extra": "59811 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 21006,
            "unit": "ns/op",
            "extra": "59329 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20622,
            "unit": "ns/op",
            "extra": "59583 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "id": "74a855787d9f6c1f70d29fe8fbf9fdb5be510dd1",
          "message": "Dev",
          "timestamp": "2024-12-09T07:50:17Z",
          "url": "https://github.com/sunerpy/requests/pull/5/commits/74a855787d9f6c1f70d29fe8fbf9fdb5be510dd1"
        },
        "date": 1733731142232,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 94765,
            "unit": "ns/op",
            "extra": "12505 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 98538,
            "unit": "ns/op",
            "extra": "12410 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 98733,
            "unit": "ns/op",
            "extra": "12212 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 98483,
            "unit": "ns/op",
            "extra": "12176 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25116,
            "unit": "ns/op",
            "extra": "46404 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 20194,
            "unit": "ns/op",
            "extra": "60032 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 20849,
            "unit": "ns/op",
            "extra": "58196 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 20421,
            "unit": "ns/op",
            "extra": "59282 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nkuzhangshn@gmail.com",
            "name": "sunerpy",
            "username": "sunerpy"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8d7d31d1a9863ff98cc35a81835e31f21e0d75b7",
          "message": "Dev (#5)\n\n* ci: 添加主分支推送测试工作流\r\n\r\n- 新增 GitHub Actions 工作流，监听主分支推送事件\r\n- 配置 Ubuntu 环境下运行测试\r\n- 包含代码检查、依赖安装、测试执行和覆盖率检查等步骤\r\n- 覆盖率低于 90% 时构建失败\r\n- 集成 Codecov 上传覆盖率报告",
          "timestamp": "2024-12-09T16:02:59+08:00",
          "tree_id": "272ab9fd9852b81b119c2f445604e1f669c1d496",
          "url": "https://github.com/sunerpy/requests/commit/8d7d31d1a9863ff98cc35a81835e31f21e0d75b7"
        },
        "date": 1733731420385,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkHTTPLibraries/NetHTTP",
            "value": 97663,
            "unit": "ns/op",
            "extra": "12210 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests",
            "value": 100931,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2",
            "value": 101007,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibraries/Requests_HTTP2_withmax",
            "value": 101062,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/NetHTTP",
            "value": 25897,
            "unit": "ns/op",
            "extra": "47550 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests",
            "value": 24963,
            "unit": "ns/op",
            "extra": "57477 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2",
            "value": 21235,
            "unit": "ns/op",
            "extra": "57405 times\n4 procs"
          },
          {
            "name": "BenchmarkHTTPLibrariesParallel/Requests_HTTP2_withmax",
            "value": 21269,
            "unit": "ns/op",
            "extra": "57969 times\n4 procs"
          }
        ]
      }
    ]
  }
}