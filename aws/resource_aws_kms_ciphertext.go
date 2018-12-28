package aws

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsKmsCiphertext() *schema.Resource {

	return &schema.Resource{
		Create: resourceAwsKmsCiphertextCreate,
		Read:   resourceAwsKmsCiphertextRead,
		Delete: resourceAwsKmsCiphertextDelete,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Sensitive: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"context": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},

			"ciphertext_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsKmsCiphertextCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).kmsconn

	d.SetId(time.Now().UTC().String())

	req := &kms.EncryptInput{
		KeyId:     aws.String(d.Get("key_id").(string)),
		Plaintext: []byte(d.Get("plaintext").(string)),
	}

	if ec := d.Get("context"); ec != nil {
		req.EncryptionContext = stringMapToPointers(ec.(map[string]interface{}))
	}

	log.Printf("[DEBUG] KMS encrypt for key: %s", d.Get("key_id").(string))
	resp, err := conn.Encrypt(req)
	if err != nil {
		return err
	}

	d.Set("ciphertext_blob", base64.StdEncoding.EncodeToString(resp.CiphertextBlob))

	return nil
}

func resourceAwsKmsCiphertextDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func resourceAwsKmsCiphertextRead(d *schema.ResourceData, meta interface{}) error {
	// If the input has changed, generate a new ciphertext_blob
	// This should never be the case since ForceNew set on all input.
	if d.HasChange("plaintext") || d.HasChange("key_id") || d.HasChange("context") {
		return resourceAwsKmsCiphertextCreate(d, meta)
	}
	return nil
}
