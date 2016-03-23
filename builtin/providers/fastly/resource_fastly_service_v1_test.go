package fastly

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	gofastly "github.com/sethvargo/go-fastly"
)

func TestResourceFastlyFlattenDomains(t *testing.T) {
	cases := []struct {
		remote []*gofastly.Domain
		local  []map[string]interface{}
	}{
		{
			remote: []*gofastly.Domain{
				&gofastly.Domain{
					Name:    "test.notexample.com",
					Comment: "not comment",
				},
			},
			local: []map[string]interface{}{
				map[string]interface{}{
					"name":    "test.notexample.com",
					"comment": "not comment",
				},
			},
		},
		{
			remote: []*gofastly.Domain{
				&gofastly.Domain{
					Name: "test.notexample.com",
				},
			},
			local: []map[string]interface{}{
				map[string]interface{}{
					"name":    "test.notexample.com",
					"comment": "",
				},
			},
		},
	}

	for _, c := range cases {
		out := flattenDomains(c.remote)
		if !reflect.DeepEqual(out, c.local) {
			t.Fatalf("Error matching:\nexpected: %#v\ngot: %#v", c.local, out)
		}
	}
}

func TestResourceFastlyFlattenBackend(t *testing.T) {
	cases := []struct {
		remote []*gofastly.Backend
		local  []map[string]interface{}
	}{
		{
			remote: []*gofastly.Backend{
				&gofastly.Backend{
					Name:                "test.notexample.com",
					Address:             "www.notexample.com",
					Port:                uint(80),
					AutoLoadbalance:     true,
					BetweenBytesTimeout: uint(10000),
					ConnectTimeout:      uint(1000),
					ErrorThreshold:      uint(0),
					FirstByteTimeout:    uint(15000),
					MaxConn:             uint(200),
					SSLCheckCert:        true,
					Weight:              uint(100),
				},
			},
			local: []map[string]interface{}{
				map[string]interface{}{
					"name":                  "test.notexample.com",
					"address":               "www.notexample.com",
					"port":                  80,
					"auto_loadbalance":      true,
					"between_bytes_timeout": 10000,
					"connect_timeout":       1000,
					"error_threshold":       0,
					"first_byte_timeout":    15000,
					"max_conn":              200,
					"ssl_check_cert":        true,
					"weight":                100,
				},
			},
		},
	}

	for _, c := range cases {
		out := flattenBackends(c.remote)
		if !reflect.DeepEqual(out, c.local) {
			t.Fatalf("Error matching:\nexpected: %#v\ngot: %#v", c.local, out)
		}
	}
}

func TestAccFastlyServiceV1_basic_updateDomain(t *testing.T) {
	var service gofastly.ServiceDetail
	name := acctest.RandString(10)
	nameUpdate := acctest.RandString(10)
	domainName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceV1Config(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV1Exists("fastly_service_v1.foo", &service),
					testAccCheckFastlyServiceV1Attributes(&service, name, domainName),
					resource.TestCheckResourceAttr(
						"fastly_service_v1.foo", "name", name),
					resource.TestCheckResourceAttr(
						"fastly_service_v1.foo", "active_version", "1"),
				),
			},

			resource.TestStep{
				Config: testAccServiceV1Config_domainUpdate(nameUpdate, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV1Exists("fastly_service_v1.foo", &service),
					testAccCheckFastlyServiceV1Attributes(&service, nameUpdate, domainName),
					resource.TestCheckResourceAttr(
						"fastly_service_v1.foo", "name", nameUpdate),
					resource.TestCheckResourceAttr(
						"fastly_service_v1.foo", "active_version", "2"),
				),
			},
		},
	})
}

func TestAccFastlyServiceV1_basic_updateBackend(t *testing.T) {
	var service gofastly.ServiceDetail
	name := acctest.RandString(10)
	// nameUpdate := acctest.RandString(10)
	backendName := acctest.RandString(3)
	backendName2 := acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceV1Config_backend(name, backendName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV1Exists("fastly_service_v1.foo", &service),
					testAccCheckFastlyServiceV1Attributes_backends(&service, name, []string{backendName}),
				),
			},

			resource.TestStep{
				Config: testAccServiceV1Config_backend_update(name, backendName, backendName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV1Exists("fastly_service_v1.foo", &service),
					testAccCheckFastlyServiceV1Attributes_backends(&service, name, []string{backendName, backendName2}),
				),
			},
		},
	})
}

func TestAccFastlyServiceV1_domain(t *testing.T) {
	var service gofastly.ServiceDetail
	name := acctest.RandString(10)
	domainName := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceV1Config(name, domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV1Exists("fastly_service_v1.foo", &service),
					testAccCheckFastlyServiceV1Attributes(&service, name, domainName),
					resource.TestCheckResourceAttr(
						"fastly_service_v1.foo", "name", name),
					resource.TestCheckResourceAttr(
						"fastly_service_v1.foo", "active_version", "1"),
				),
			},
		},
	})
}

func testAccCheckServiceV1Exists(n string, service *gofastly.ServiceDetail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Service ID is set")
		}

		conn := testAccProvider.Meta().(*FastlyClient).conn
		latest, err := conn.GetServiceDetails(&gofastly.GetServiceInput{
			ID: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		*service = *latest

		return nil
	}
}

func testAccCheckFastlyServiceV1Attributes(service *gofastly.ServiceDetail, name, domain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if service.Name != name {
			return fmt.Errorf("Bad name, expected (%s), got (%s)", name, service.Name)
		}

		log.Printf("\n---\nDEBUG\n---\nService:\n%#v\n\n---\n", service)

		for _, v := range service.Versions {
			log.Printf("\n\tversion: %#v\n\n", v)
		}

		return nil
	}
}

func testAccCheckFastlyServiceV1Attributes_backends(service *gofastly.ServiceDetail, name string, backends []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if service.Name != name {
			return fmt.Errorf("Bad name, expected (%s), got (%s)", name, service.Name)
		}

		log.Printf("\n---\nDEBUG\n---\nService:\n%#v\n\n---\n", service)

		conn := testAccProvider.Meta().(*FastlyClient).conn
		l, err := conn.ListBackends(&gofastly.ListBackendsInput{
			Service: service.ID,
			Version: service.ActiveVersion.Number,
		})

		log.Printf("\n---\nDEBUG: l : %#v\n---\n", l)

		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckServiceV1Destroy(s *terraform.State) error {
	// conn := testAccProvider.Meta().(*FastlyClient).conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "fastly_service_v1" {
			continue
		}

		conn := testAccProvider.Meta().(*FastlyClient).conn
		l, err := conn.ListServices(&gofastly.ListServicesInput{})
		if err != nil {
			return fmt.Errorf("[WARN] Error listing servcies when deleting Fastly Service (%s): %s", rs.Primary.ID, err)
		}

		for _, s := range l {
			if s.ID == rs.Primary.ID {
				// service still found
				return fmt.Errorf("[WARN] Tried deleting Service (%s), but was still found", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccServiceV1Config(name, domain string) string {
	return fmt.Sprintf(`
resource "fastly_service_v1" "foo" {
  name = "%s"

  domain {
    name    = "%s.notadomain.com"
    comment = "tf-testing-domain"
  }

  backend {
    address = "aws.amazon.com"
    name    = "amazon docs"
  }

  force_destroy = true
}`, name, domain)
}

func testAccServiceV1Config_domainUpdate(name, domain string) string {
	return fmt.Sprintf(`
resource "fastly_service_v1" "foo" {
  name = "%s"

  domain {
    name    = "%s.notadomain.com"
    comment = "tf-testing-domain"
  }

  domain {
    name    = "%s.notanotherdomain.com"
    comment = "tf-testing-other-domain"
  }

  backend {
    address = "aws.amazon.com"
    name    = "amazon docs"
  }

  force_destroy = true
}`, name, domain, domain)
}

func testAccServiceV1Config_backend(name, backend string) string {
	return fmt.Sprintf(`
resource "fastly_service_v1" "foo" {
  name = "%s"

  domain {
    name    = "test.notadomain.com"
    comment = "tf-testing-domain"
  }

  backend {
    address = "%s.aws.amazon.com"
    name    = "tf -test backend"
  }

  force_destroy = true
}`, name, backend)
}

func testAccServiceV1Config_backend_update(name, backend, backend2 string) string {
	return fmt.Sprintf(`
resource "fastly_service_v1" "foo" {
  name = "%s"

  domain {
    name    = "test.notadomain.com"
    comment = "tf-testing-domain"
  }

  backend {
    address = "%s.aws.amazon.com"
    name    = "tf-test-backend"
  }

  backend {
    address = "%s.aws.amazon.com"
    name    = "tf-test-backend-other"
  }

  force_destroy = true
}`, name, backend, backend2)
}
